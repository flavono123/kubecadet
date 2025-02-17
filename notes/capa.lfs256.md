---
id: hrrz5jmtnbwdzptwacyigdx
title: Lfs256
desc: ''
updated: 1739803252086
created: 1738998769740
---

- [Milestones](#milestones)
- [01. Course Introduction](#01-course-introduction)
- [02. Introduction to Argo](#02-introduction-to-argo)
- [03. Argo CD](#03-argo-cd)
  - [Vocab](#vocab)
  - [Components](#components)
  - [Reconciliation](#reconciliation)
  - [Syncronization Principles](#syncronization-principles)
    - [Resource Hooks(Phase)](#resource-hooksphase)
    - [Sync Wave](#sync-wave)
    - [Orders](#orders)
  - [Objects \& Resources](#objects--resources)
    - [Application](#application)
    - [AppProject](#appproject)
    - [Repository credentials](#repository-credentials)
    - [Cluster credentials](#cluster-credentials)
  - [Plugins](#plugins)
  - [Secure](#secure)
    - [RBAC](#rbac)
- [04. Argo Workflows](#04-argo-workflows)
  - [Objectives Keywords](#objectives-keywords)
  - [Core Concepts](#core-concepts)
    - [Workflow](#workflow)
    - [Outputs(Artifacts)](#outputsartifacts)
    - [Workflow Templates](#workflow-templates)
    - [Cluster Workflow Templates(v2.8+)](#cluster-workflow-templatesv28)
    - [Cron Workflows(v2.5+)](#cron-workflowsv25)
  - [Architecture](#architecture)
    - [Containers of Pod(a step or dag task)](#containers-of-poda-step-or-dag-task)
  - [Use cases](#use-cases)


## Milestones

- [x] day 1 01. Course Introduction & 02. Introduction to Argo - 2/10(mon)
- [x] day 2 03. Argo CD  - 2/15(sat)
- [x] day 3 04. Argo Workflows 1 - 2/16(sun)
- [ ] day 4 04. Argo Workflows 2 - 2/17(mon)
- [ ] day 5 05. Argo Rollouts - 2/18(tue)
- [ ] day 6 06. Argo Events - 2/19(wed)

## 01. Course Introduction

- volume: 16-20hours; 4h per day?

## 02. Introduction to Argo

GitOps Core Elements:

- Declarative configuration
- Immutable storage
- Automation(post-commit works)
- Software agents
- Closed loop

## 03. Argo CD

### Vocab

- Application: CRD, a collection of k8s resources
- Application source types: Helm, Kustomize, ...
- State: Target state, Live state
- Status
  - Sync: Live state == Target state
  - Sync operation: failed, succeeded
  - Health: of application
- Actions
  - Refresh: identify the difference between latest code in git repo
  - Sync: apply as the target state

### Components

- controller: watch resources' spec
- api server <- user, ci/cd
  - mgmt apps, status updates
  - trigger operations on apps
  - handles git repos for version control
  - connect with k8s clusters
  - auth, sso, rbac
  - central communication hub; web ui, cli, argo events, ...
- repository server -> git
  - argo req manifest args: repo url, revision, app path, template rel info(e.g. helm values)
- application controller -> k8s objects
  - detect discrepancies

### Reconciliation

- monitors yaml manifests from git repo
- observes k8s cluster/objects, `kubectl apply` when a disparity is identified(if auto-sync is enabled)

### Syncronization Principles

#### Resource Hooks(Phase)

- PreSync
- Sync; after PreSync, **during Sync application**(e.g. complex rollout strategy such as blue/green, canary, ... )
- PostSync; after Sync when all resource is in Healthy state
- Skip
- SyncFail; when Sync failed(e.g. clean up)
- PostDelete; after all Applications are deleted(v2.1.0+)
- use a Job specific hook annotated(e.g. `argocd.argoproj.io/hook: PreSync`)
- ref. [resource hooks doc](https://argo-cd.readthedocs.io/en/stable/user-guide/resource_hooks/#resource-hooks)

#### Sync Wave

- for splitting and ordering manifest syncs
- negative to positive values; default is wave 0
- annotate to resource manifests such as `argocd.argoproj.io/sync-wave: "-1"`
- for each wave, ARGOCD_SYNC_WAVE_DELAY is applied; default is 2s
- ref. [sync waves doc](https://argo-cd.readthedocs.io/en/stable/user-guide/sync-waves/#how-do-i-configure-waves)

#### Orders

1. hooks
2. sync waves
3. [Kind priority](https://github.com/argoproj/gitops-engine/blob/bc9ce5764fa306f58cf59199a94f6c968c775a2d/pkg/sync/sync_tasks.go#L27-L66)
4. [Name(alphabetical)](https://github.com/argoproj/gitops-engine/blob/65db274b8d73302f131768571ff1bb9383f476af/pkg/sync/sync_tasks.go#L80-L110)

### Objects & Resources

#### Application

[Application CRD](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/#applications)

#### AppProject

[AppProject CRD](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/#projects), grouping applications

#### Repository credentials

Secret

- labeled `argocd.argoproj.io/secret-type: repository`
- data: `url`, auth methods(`username`, `password`, ...)

#### Cluster credentials

Secret:

- labeled `argocd.argoproj.io/secret-type: cluster`
- data: `config`, `name`, `server`

### Plugins

e.g. Notification (it seems one of core components for now)

- configuring with ConfigMap
- ref. [notification doc](https://argo-cd.readthedocs.io/en/stable/operator-manual/notifications/)
- ref. [plugin doc](https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/)
- plugins:
  - [image updater](https://argocd-image-updater.readthedocs.io/en/stable/); automates updating container images
  - [autopilot](https://argocd-autopilot.readthedocs.io/en/stable/); bootstrap a new argocd to a cluster
  - [interlace](https://gisthub.com/argoproj-labs/argocd-interlace)

### Secure

#### RBAC

> define roles in `argocd-rbac-cm`, streamlines to k8s rbac resources, bindings and principals

ref. [security doc](https://argo-cd.readthedocs.io/en/stable/operator-manual/security/)

## 04. Argo Workflows

### Objectives Keywords

- metadata, spec, entrypoint, templates
- controller, ui
- scheduling, execution

### Core Concepts

#### Workflow

- define and store a state of a workflow
- specified for automated steps involved in the deployment, testing and promotion of apps
- two main parts of spec:
  - `entrypoint`: the name of the template starting point of the workflow execution
  - `templates`: steps(tasks) to be executed, 9 types are.
    - [definitions](https://argo-workflows.readthedocs.io/en/latest/workflow-concepts/#core-concepts)
      - [container](https://argo-workflows.readthedocs.io/en/latest/workflow-concepts/#container)
      - [container set](https://argo-workflows.readthedocs.io/en/latest/container-set-template/#containerset-template)
      - [resource](https://argo-workflows.readthedocs.io/en/latest/walk-through/kubernetes-resources/) - k8s objects

    ```yaml
    templates:
    - name: resource
      resource:
        action: create            # kubectl action (e.g. create, delete, apply, patch)
        # label selection syntax and can be applied against any field of the resource (not just labels)
        # Multiple AND conditions can be represented by comma delimited expressions.
        # For more details: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
        # successCondition: status.succeeded > 0
        # failureCondition: status.failed > 3
        manifest: |               #put your kubernetes spec here
          apiVersion: batch/v1
          kind: Job
          ...
    ```

      - [script](https://argo-workflows.readthedocs.io/en/latest/http-template/)
      - [suspend](https://argo-workflows.readthedocs.io/en/latest/walk-through/suspending/)
      - [plugin](**https**://argo-workflows.readthedocs.io/en/latest/http-template/)
      - + [data sources, transformations](https://argo-workflows.readthedocs.io/en/latest/data-sourcing-and-transformation/)
      - + [http](https://argo-workflows.readthedocs.io/en/latest/http-template/)
        - use the argo agent, dedicated controller to requests
    - [invocators](https://argo-workflows.readthedocs.io/en/latest/workflow-concepts/#template-invocators)
      - dag
      - steps

#### Outputs(Artifacts)

- define by `outputs` of a template, reference that in another step using templating [expressions](https://argo-workflows.readthedocs.io/en/latest/variables/#expression)([workflow inputs](https://argo-workflows.readthedocs.io/en/latest/workflow-inputs/#workflow-inputs))

#### [Workflow Templates](https://argo-workflows.readthedocs.io/en/latest/workflow-templates/)

- reusing workflow(`templates` in other words, tasks) directly or referencing
- the fields of WorkflowSpec, except `priority`, are compatible with WorkflowTemplateSpec (v2.7+)
  - `entrypoint` is not allowed v2.4 ~ v2.6
- `workflowMetata` is for adding labels/annotations to a templated workflow
- paramters are for passing parameters to a templated workflow by templating expressions
  - a global param is by `arguments.parameters` of a workflow
  - a local param is by `templates[].inputs`
- referencing by `templates.{invocator}.templateRef` of a workflow
- creating a workflow with `workflowTemplateRef` would be merged with a referenced workflow template (v2.9+)

#### [Cluster Workflow Templates(v2.8+)](https://argo-workflows.readthedocs.io/en/latest/cluster-workflow-templates/#cluster-workflow-templates)

- a cluster scoped workflow template, such as ClusterRole
- for referencing, `templateRef.clusterScope: true` is required for (specname)
- for creating a workflow, `workflowTemplateRef.clusterScope: true` is required

#### [Cron Workflows(v2.5+)](https://argo-workflows.readthedocs.io/en/latest/cron-workflows/#cron-workflows)

- running workflows on schedule by cron expression, mimic of CronJob
- `workflowSpec` and `workflowMetadata` are allowed, different from a cronjob using `jobTemplate`.
- [options](https://argo-workflows.readthedocs.io/en/latest/cron-workflows/#cronworkflow-options)
  - difference with CronJob:
    - `schedules` is for a list of cron expressions
    - supporting `timezone`
  - unique features:
    - `stopStrategy.expression` (nil) - if expression is evaluated to false, stop the workflow
    - `when` (None) - additional condition for running a workflow on cron schedules

### [Architecture](https://argo-workflows.readthedocs.io/en/latest/architecture/)

- argo server: the api server for workflow submission, monitoring and management
- workflow controller: manage lifecycle of workflows, watch CRs

#### Containers of Pod(a step or dag task)

- init: InitContainer, fetching artifcats, parameters and making them available for main container
- main: runs the Image that the user indicated, where the argoexec utility is volume mounted and serves as the main command which calls the configured Command as a sub-process
- wait: cleanup, saving off artifacts and parameters

### Use cases

- data processing
- ml projects
- ci/cd
- batch processing

