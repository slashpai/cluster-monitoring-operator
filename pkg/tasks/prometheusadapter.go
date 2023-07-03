package tasks

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
)

type PrometheusAdapterTask struct {
	client    *client.Client
	ctx       context.Context
	factory   *manifests.Factory
	config    *manifests.Config
	namespace string
}

func NewPrometheusAdapterTask(ctx context.Context, namespace string, client *client.Client, factory *manifests.Factory, config *manifests.Config) *PrometheusAdapterTask {
	return &PrometheusAdapterTask{
		client:    client,
		factory:   factory,
		config:    config,
		namespace: namespace,
		ctx:       ctx,
	}
}

func (t *PrometheusAdapterTask) Run(ctx context.Context) error {
	if t.config.ClusterMonitoringConfiguration.MetricsServerConfig != nil {
		if !*t.config.ClusterMonitoringConfiguration.MetricsServerConfig.Enabled {
			t.create(ctx)
		}
	}

	return nil
}

func (t *PrometheusAdapterTask) create(ctx context.Context) error {
	{
		cr, err := t.factory.PrometheusAdapterClusterRole()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ClusterRole failed")
		}

		err = t.client.CreateOrUpdateClusterRole(ctx, cr)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ClusterRole failed")
		}
	}
	{
		cr, err := t.factory.PrometheusAdapterClusterRoleServerResources()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ClusterRole for server resources failed")
		}

		err = t.client.CreateOrUpdateClusterRole(ctx, cr)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ClusterRole for server resources failed")
		}
	}
	{
		cr, err := t.factory.PrometheusAdapterClusterRoleAggregatedMetricsReader()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ClusterRole aggregating resource metrics read permissions failed")
		}

		err = t.client.CreateOrUpdateClusterRole(ctx, cr)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ClusterRole aggregating resource metrics read permissions failed")
		}
	}
	{
		crb, err := t.factory.PrometheusAdapterClusterRoleBinding()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ClusterRoleBinding failed")
		}

		err = t.client.CreateOrUpdateClusterRoleBinding(ctx, crb)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ClusterRoleBinding failed")
		}
	}
	{
		crb, err := t.factory.PrometheusAdapterClusterRoleBindingDelegator()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ClusterRoleBinding for delegator failed")
		}

		err = t.client.CreateOrUpdateClusterRoleBinding(ctx, crb)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ClusterRoleBinding for delegator failed")
		}
	}
	{
		crb, err := t.factory.PrometheusAdapterClusterRoleBindingView()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ClusterRoleBinding for view failed")
		}

		err = t.client.CreateOrUpdateClusterRoleBinding(ctx, crb)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ClusterRoleBinding for view failed")
		}
	}
	{
		rb, err := t.factory.PrometheusAdapterRoleBindingAuthReader()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter RoleBinding for auth-reader failed")
		}

		err = t.client.CreateOrUpdateRoleBinding(ctx, rb)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter RoleBinding for auth-reader failed")
		}
	}
	{
		sa, err := t.factory.PrometheusAdapterServiceAccount()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ServiceAccount failed")
		}

		err = t.client.CreateOrUpdateServiceAccount(ctx, sa)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ServiceAccount failed")
		}
	}
	{
		cm, err := t.factory.PrometheusAdapterConfigMapAuditPolicy()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter AuditPolicy ConfigMap failed")
		}

		err = t.client.CreateOrUpdateConfigMap(ctx, cm)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter AuditPolicy ConfigMap failed")
		}
	}
	{
		cm, err := t.factory.PrometheusAdapterConfigMapPrometheus()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ConfigMap for Prometheus failed")
		}

		err = t.client.CreateOrUpdateConfigMap(ctx, cm)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter ConfigMap for Prometheus failed")
		}
	}
	{
		s, err := t.factory.PrometheusAdapterService()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter Service failed")
		}

		err = t.client.CreateOrUpdateService(ctx, s)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter Service failed")
		}
	}
	// Intermediate variable to hold on to the config map name that the
	// prometheus-adapter deployment should target.
	var cmName string
	{
		cmD, err := t.factory.PrometheusAdapterConfigMapDedicated()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ConfigMap for dedicated ServiceMonitors failed")
		}
		cm, err := t.factory.PrometheusAdapterConfigMap()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ConfigMap failed")
		}
		if t.config.ClusterMonitoringConfiguration.K8sPrometheusAdapter.DedicatedServiceMonitors.Enabled {
			err = t.client.CreateOrUpdateConfigMap(ctx, cmD)
			if err != nil {
				return errors.Wrap(err, "reconciling PrometheusAdapter ConfigMap for dedicated ServiceMonitors failed")
			}
			err = t.client.DeleteConfigMap(ctx, cm)
			if err != nil {
				return errors.Wrap(err, "deleting PrometheusAdapter ConfigMap failed")
			}
			cmName = cmD.Name
		} else {
			err = t.client.CreateOrUpdateConfigMap(ctx, cm)
			if err != nil {
				return errors.Wrap(err, "reconciling PrometheusAdapter ConfigMap failed")
			}
			err = t.client.DeleteConfigMap(ctx, cmD)
			if err != nil {
				return errors.Wrap(err, "deleting PrometheusAdapter ConfigMap for dedicated ServiceMonitors failed")
			}
			cmName = cm.Name
		}

		tlsSecret, err := t.client.WaitForSecretByNsName(ctx, types.NamespacedName{Namespace: t.namespace, Name: "prometheus-adapter-tls"})
		if err != nil {
			return errors.Wrap(err, "failed to wait for prometheus-adapter-tls secret")
		}

		apiAuthConfigmap, err := t.client.WaitForConfigMapByNsName(ctx, types.NamespacedName{Namespace: "kube-system", Name: "extension-apiserver-authentication"})
		if err != nil {
			return errors.Wrap(err, "failed to wait for kube-system/extension-apiserver-authentication configmap")
		}

		secret, err := t.factory.PrometheusAdapterSecret(tlsSecret, apiAuthConfigmap)
		if err != nil {
			return errors.Wrap(err, "failed to create prometheus adapter secret")
		}

		err = t.deleteOldPrometheusAdapterSecrets(secret.Labels["monitoring.openshift.io/hash"])
		if err != nil {
			return errors.Wrap(err, "deleting old prometheus adapter secrets failed")
		}

		err = t.client.CreateOrUpdateSecret(ctx, secret)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter Secret failed")
		}

		dep, err := t.factory.PrometheusAdapterDeployment(secret.Name, apiAuthConfigmap.Data, cmName)
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter Deployment failed")
		}

		err = t.client.CreateOrUpdateDeployment(ctx, dep)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter Deployment failed")
		}
	}
	{
		pdb, err := t.factory.PrometheusAdapterPodDisruptionBudget()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter PodDisruptionBudget failed")
		}

		if pdb != nil {
			err = t.client.CreateOrUpdatePodDisruptionBudget(ctx, pdb)
			if err != nil {
				return errors.Wrap(err, "reconciling PrometheusAdapter PodDisruptionBudget failed")
			}
		}
	}
	{
		sms, err := t.factory.PrometheusAdapterServiceMonitors()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter ServiceMonitors failed")
		}

		for _, sm := range sms {
			err = t.client.CreateOrUpdateServiceMonitor(ctx, sm)
			if err != nil {
				return errors.Wrapf(err, "reconciling %s/%s ServiceMonitor failed", sm.Namespace, sm.Name)
			}
		}
	}
	{
		api, err := t.factory.PrometheusAdapterAPIService()
		if err != nil {
			return errors.Wrap(err, "initializing PrometheusAdapter APIService failed")
		}

		err = t.client.CreateOrUpdateAPIService(ctx, api)
		if err != nil {
			return errors.Wrap(err, "reconciling PrometheusAdapter APIService failed")
		}
	}
	return nil
}

func (t *PrometheusAdapterTask) deleteOldPrometheusAdapterSecrets(newHash string) error {
	secrets, err := t.client.KubernetesInterface().CoreV1().Secrets(t.namespace).List(t.ctx, metav1.ListOptions{
		LabelSelector: "monitoring.openshift.io/name=prometheus-adapter,monitoring.openshift.io/hash!=" + newHash,
	})

	if err != nil {
		return errors.Wrap(err, "error listing prometheus adapter secrets")
	}

	for i := range secrets.Items {
		err := t.client.KubernetesInterface().CoreV1().Secrets(t.namespace).Delete(t.ctx, secrets.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrapf(err, "error deleting secret: %s", secrets.Items[i].Name)
		}
	}

	return nil
}
