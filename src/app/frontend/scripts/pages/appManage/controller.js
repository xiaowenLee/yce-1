/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope', 'appManageService', '$sessionStorage', function($scope,appManageService, $sessionStorage){

            $scope.param = {"orgId": $sessionStorage.orgId, "userId": $sessionStorage.userId}

            appManageService.getAppList($scope.param,function(data){
                // if (data.code == 0) {
                //   $scope.appList = data.data;
                // }

                // Tobedeleted
                var appList = {
                                             "dcId": "bangongwang",
                                             "dcName": "世纪互联",
                                             "podList": {
                                                  "kind": "PodList",
                                                    "apiVersion": "v1",
                                                    "metadata": {
                                                      "selfLink": "/api/v1/namespaces/ops/pods",
                                                      "resourceVersion": "16266621"
                                                    },
                                                    "items": [
                                                      {
                                                        "metadata": {
                                                          "name": "nginx-test-1252813378-39wjk",
                                                          "generateName": "nginx-test-1252813378-",
                                                          "namespace": "ops",
                                                          "selfLink": "/api/v1/namespaces/ops/pods/nginx-test-1252813378-39wjk",
                                                          "uid": "217b98e7-64ee-11e6-b957-44a84240716a",
                                                          "resourceVersion": "16258268",
                                                          "creationTimestamp": "2016-08-18T02:47:40Z",
                                                          "labels": {
                                                            "name": "spec-template-metadata-labels",
                                                            "pod-template-hash": "1252813378"
                                                          },
                                                          "annotations": {
                                                            "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"ops\",\"name\":\"nginx-test-1252813378\",\"uid\":\"2179503e-64ee-11e6-b957-44a84240716a\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"16258237\"}}"
                                                          }
                                                        },
                                                        "spec": {
                                                          "volumes": [
                                                            {
                                                              "name": "default-token-jr2e0",
                                                              "secret": {
                                                                "secretName": "default-token-jr2e0"
                                                              }
                                                            }
                                                          ],
                                                          "containers": [
                                                            {
                                                              "name": "nginx-test",
                                                              "image": "nginx:1.7.9",
                                                              "resources": {},
                                                              "volumeMounts": [
                                                                {
                                                                  "name": "default-token-jr2e0",
                                                                  "readOnly": true,
                                                                  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
                                                                }
                                                              ],
                                                              "terminationMessagePath": "/dev/termination-log",
                                                              "imagePullPolicy": "IfNotPresent"
                                                            }
                                                          ],
                                                          "restartPolicy": "Always",
                                                          "terminationGracePeriodSeconds": 30,
                                                          "dnsPolicy": "ClusterFirst",
                                                          "serviceAccountName": "default",
                                                          "serviceAccount": "default",
                                                          "nodeName": "172.21.1.21",
                                                          "securityContext": {}
                                                        },
                                                        "status": {
                                                          "phase": "Running",
                                                          "conditions": [
                                                            {
                                                              "type": "Ready",
                                                              "status": "True",
                                                              "lastProbeTime": null,
                                                              "lastTransitionTime": "2016-08-18T02:47:41Z"
                                                            }
                                                          ],
                                                          "hostIP": "172.21.1.21",
                                                          "podIP": "10.0.62.3",
                                                          "startTime": "2016-08-18T02:47:40Z",
                                                          "containerStatuses": [
                                                            {
                                                              "name": "nginx-test",
                                                              "state": {
                                                                "running": {
                                                                  "startedAt": "2016-08-18T02:47:41Z"
                                                                }
                                                              },
                                                              "lastState": {},
                                                              "ready": true,
                                                              "restartCount": 0,
                                                              "image": "nginx:1.7.9",
                                                              "imageID": "docker://sha256:84581e99d807a703c9c03bd1a31cd9621815155ac72a7365fd02311264512656",
                                                              "containerID": "docker://0ac0b1d2e4bc085a1655049e7c251c056bffe05ec26b60523bd34fd590bcc472"
                                                            }
                                                          ]
                                                        }
                                                      }
                                                      ]
                                             }
                                 };

            });

            $scope.appList = appList;
            console.log(appList);
            console.dir($scope.appList)

        }];


        var controllers = [
            {module: 'appManage', name: 'appManageController', ctrl: ctrl}
        ];

        return controllers;
    }
);