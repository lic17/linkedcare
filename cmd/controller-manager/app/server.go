/*

 Copyright 2019 The Linkedcare Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

package app

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	networkingapis "istio.io/client-go/pkg/apis/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"linkedcare.io/linkedcare/cmd/controller-manager/app/options"
	apis "linkedcare.io/linkedcare/pkg/apis/servicemesh/v1alpha1"
	"linkedcare.io/linkedcare/pkg/client/clientset/versioned/scheme"
	"linkedcare.io/linkedcare/pkg/controller"
	"linkedcare.io/linkedcare/pkg/controller/client"
	appapis "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func NewControllerManagerCommand() *cobra.Command {
	s := options.NewLinkedcareControllerManagerOptions()

	cmd := &cobra.Command{
		Use:  "controller-manager",
		Long: `Linkedcare controller manager is a daemon that`,
		Run: func(cmd *cobra.Command, args []string) {

			if err := Run(s, signals.SetupSignalHandler()); err != nil {
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	return cmd
}

func Run(s *options.LinkedcareControllerManagerOptions, stopCh <-chan struct{}) error {
	clientset := client.CreateClientSet(s.Kubeconfig, stopCh)

	config := clientset.Config

	run := func(ctx context.Context) {
		klog.V(0).Info("setting up manager")
		mgr, err := manager.New(config, manager.Options{})
		if err != nil {
			klog.Fatalf("unable to set up overall controller manager: %v", err)
		}

		klog.V(0).Info("setting up scheme")
		if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
			klog.Fatalf("unable add APIs to scheme: %v", err)
		}

		klog.V(0).Info("setting up istio scheme")
		if err := networkingapis.AddToScheme(mgr.GetScheme()); err != nil {
			klog.Fatalf("unable add istio APIs to scheme: %v", err)
		}

		klog.V(0).Info("setting up application scheme")
		if err := appapis.AddToScheme(mgr.GetScheme()); err != nil {
			klog.Fatalf("unable add application APIs to scheme: %v", err)
		}

		klog.V(0).Info("Setting up controllers")
		if err := controller.AddToManager(mgr); err != nil {
			klog.Fatalf("unable to register controllers to the manager: %v", err)
		}

		if err := AddControllers(mgr, config, stopCh); err != nil {
			klog.Fatalf("unable to register controllers to the manager: %v", err)
		}

		klog.V(0).Info("Starting the Cmd.")
		if err := mgr.Start(stopCh); err != nil {
			klog.Fatalf("unable to run the manager: %v", err)
		}

		select {}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stopCh
		cancel()
	}()

	id, err := os.Hostname()
	if err != nil {
		return err
	}

	// add a uniquifier so that two processes on the same host don't accidentally both become active
	id = id + "_" + string(uuid.NewUUID())

	// TODO: change lockType to lease
	// once we finished moving to Kubernetes v1.16+, we
	// change lockType to lease
	lock, err := resourcelock.New(resourcelock.LeasesResourceLock,
		"linkedcare-system",
		"lk-controller-manager",
		clientset.K8sClient.CoreV1(),
		clientset.K8sClient.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: id,
			EventRecorder: record.NewBroadcaster().NewRecorder(scheme.Scheme, v1.EventSource{
				Component: "ks-controller-manager",
			}),
		})

	if err != nil {
		klog.Fatalf("error creating lock: %v", err)
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: s.LeaderElection.LeaseDuration,
		RenewDeadline: s.LeaderElection.RenewDeadline,
		RetryPeriod:   s.LeaderElection.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				klog.Errorf("leadership lost")
				os.Exit(0)
			},
		},
	})

	return nil
}
