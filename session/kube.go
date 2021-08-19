package session

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"

	"github.com/evilsocket/islazy/fs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/evilsocket/islazy/async"
)

type KubeSession struct {
	sync.Mutex
	namespace     string
	podName       string
	containerName string
	config        *restclient.Config
	client        kubernetes.Interface
	timeouts      core.Timeouts
}

func NewKube(sh models.Shell, timeouts core.Timeouts) (error, Session) {
	tokenFile, err := fs.Expand(sh.Identity.KeyFile)
	if err != nil {
		return err, nil
	}

	config := &restclient.Config{
		Host:            sh.Host,
		BearerTokenFile: tokenFile,
		TLSClientConfig: restclient.TLSClientConfig{
			Insecure: sh.Insecure,
		},
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err, nil
	}

	return nil, &KubeSession{
		namespace:     sh.Namespace,
		podName:       sh.Pod,
		containerName: sh.Container,
		config:        config,
		client:        client,
		timeouts:      timeouts,
	}
}

func (k *KubeSession) Type() string {
	return "kube"
}

func (k *KubeSession) Exec(cmd string) ([]byte, error) {
	k.Lock()
	defer k.Unlock()

	obj, err := async.WithTimeout(k.timeouts.RW(), func() interface{} {
		var stdout, stderr bytes.Buffer

		req := k.client.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(k.podName).
			Namespace(k.namespace).
			SubResource("exec").
			Param("container", k.containerName)

		req.VersionedParams(&v1.PodExecOptions{
			Container: k.containerName,
			Command:   []string{"sh", "-c", cmd},
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(k.config, "POST", req.URL())
		if err != nil {
			return err
		}

		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		})
		if err != nil {
			return err
		}

		res := cmdResult{
			out: stdout.Bytes(),
		}
		out_err := strings.TrimSpace(stderr.String())
		if len(out_err) != 0 {
			res.err = fmt.Errorf("%s", out_err)
		}

		return res
	})
	if err != nil {
		return nil, err
	}

	res := obj.(cmdResult)
	return res.out, res.err
}

func (k *KubeSession) Close() {

}
