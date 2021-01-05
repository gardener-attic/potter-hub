package testing

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	hubv1 "github.com/gardener/potter-hub/pkg/external/hubcontroller/api/v1"
)

func CreateClusterBom(clusterBomName, overallState string) *hubv1.ClusterBom {
	testClusterBom := &hubv1.ClusterBom{
		ObjectMeta: v1.ObjectMeta{
			Name: clusterBomName,
		},

		Spec: hubv1.ClusterBomSpec{
			SecretRef: "testsecret01",
		},

		Status: hubv1.ClusterBomStatus{
			OverallState: overallState,
		},
	}
	return testClusterBom
}

func AddApplicationConfig(clusterBom *hubv1.ClusterBom, id string) {
	clusterBom.Spec.ApplicationConfigs = append(clusterBom.Spec.ApplicationConfigs, hubv1.ApplicationConfig{
		ID: id,
	})
}

func AddApplicationStatus(clusterBom *hubv1.ClusterBom, id, state string, currentOp *hubv1.CurrentOperation, lastOp *hubv1.LastOperation) {
	clusterBom.Status.ApplicationStates = append(clusterBom.Status.ApplicationStates, hubv1.ApplicationState{
		ID:    id,
		State: state,
		DetailedState: hubv1.DetailedState{
			CurrentOperation: *currentOp,
			LastOperation:    *lastOp,
		},
	})
}

func CreateDeploymentConfig(clusterBomName, id string, currentOp *hubv1.CurrentOperation, lastOp *hubv1.LastOperation, readiness *hubv1.Readiness) *hubv1.HubDeploymentConfig {
	return &hubv1.HubDeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name: clusterBomName + "-" + id,
			Labels: map[string]string{
				"hub.kubernetes.sap.com/bom-name": clusterBomName,
			},
		},
		Spec: hubv1.HubDeploymentConfigSpec{
			DeploymentConfig: hubv1.DeploymentConfig{
				ID: id,
			},
			CurrentOperation: *currentOp,
		},
		Status: hubv1.HubDeploymentConfigStatus{
			LastOperation: *lastOp,
			Readiness:     readiness,
		},
	}
}

func CurrentOp(operation string, number int32) hubv1.CurrentOperation {
	return hubv1.CurrentOperation{
		Operation: operation,
		Number:    number,
	}
}

func LastOp(operation string, number, successNumber int32, state string, numOfTries int32) hubv1.LastOperation {
	return hubv1.LastOperation{
		Operation:     operation,
		Number:        number,
		SuccessNumber: successNumber,
		State:         state,
		NumberOfTries: numOfTries,
	}
}

func CreateSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
	}
}
