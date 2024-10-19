import React from 'react';
import HubspotAuthorizations from "./views/pages/integrations/HubspotAuthorizations";
import SalesforceAuthorizations from "./views/pages/integrations/SalesforceAuthorizations";

import {PIPELINE_TYPE_GENERIC, PIPELINE_TYPE_METRIC} from "./services/pipeline.service";

const Dashboard = React.lazy(() => import('./views/dashboard/Dashboard'));
const Pipelines = React.lazy(() => import('./views/pages/pipelines/Pipelines'));
const Integrations = React.lazy(() => import('./views/pages/integrations/Integrations'));
const SlackAuthorizations = React.lazy(() => import('./views/pages/integrations/SlackAuthorizations'));
const JiraAuthorizations = React.lazy(() => import('./views/pages/integrations/JiraAuthorizations'));
const Settings = React.lazy(() => import('./views/pages/settings/Settings'));
const Users = React.lazy(() => import('./views/pages/users/Users'));
const Metrics = React.lazy(() => import('./views/pages/metrics/Metrics'));
const EnvVariables = React.lazy(() => import('./views/pages/env-variables/EnvVariables'));
const APIClients = React.lazy(() => import('./views/pages/api-clients/APIClients'));
const SlowTasks = React.lazy(() => import('./views/pages/profiling/SlowTasks'));

const routes = [
  { path: '/', name: 'Home', element: <Dashboard/> },
  { path: '/dashboard/:itemUuid/:dateFrom/:dateTo', name: 'Dashboard', element: <Dashboard/> },
  { path: '/dashboard', name: 'Dashboard', element: <Dashboard/> },
  { path: '/pipelines/folder/:folderUuid/:page/:itemUuid/:dateFrom/:dateTo', name: 'Pipelines', element: <Pipelines type={PIPELINE_TYPE_GENERIC}/> },
  { path: '/pipelines/folder/:folderUuid/:page/:itemUuid', name: 'Pipelines', element: <Pipelines type={PIPELINE_TYPE_GENERIC}/> },
  { path: '/pipelines/folder/:folderUuid', name: 'Pipelines', element: <Pipelines type={PIPELINE_TYPE_GENERIC}/> },
  { path: '/pipelines/folder', name: 'Pipelines', element: <Pipelines type={PIPELINE_TYPE_GENERIC}/> },
  { path: '/pipelines', name: 'Pipelines', element: <Pipelines type={PIPELINE_TYPE_GENERIC}/> },
  { path: '/metrics/folder/:folderUuid/:page/:itemUuid/:dateFrom/:dateTo', name: 'Metrics', element: <Pipelines type={PIPELINE_TYPE_METRIC}/> },
  { path: '/metrics/folder/:folderUuid/:page/:itemUuid', name: 'Metrics', element: <Pipelines type={PIPELINE_TYPE_METRIC}/> },
  { path: '/metrics/folder/:folderUuid', name: 'Metrics', element: <Pipelines type={PIPELINE_TYPE_METRIC}/> },
  { path: '/metrics/folder', name: 'Metrics', element: <Pipelines type={PIPELINE_TYPE_METRIC}/> },
  { path: '/metrics', name: 'Metrics', element: <Pipelines type={PIPELINE_TYPE_METRIC}/> },
  { path: '/integrations/:page/:itemUuid', name: 'Integrations', element: <Integrations/> },
  { path: '/integrations', name: 'Integrations', element: <Integrations/> },
  { path: '/slack-authorizations/:itemUuid', name: 'Slack Authorization', element: <SlackAuthorizations/> },
  { path: '/slack-authorizations', name: 'Slack Authorizations', element: <SlackAuthorizations/> },
  { path: '/jira-authorizations/:itemUuid', name: 'Jira Authorization', element: <JiraAuthorizations/> },
  { path: '/jira-authorizations', name: 'Jira Authorizations', element: <JiraAuthorizations/> },
  { path: '/hubspot-authorizations/:itemUuid', name: 'Hubspot Authorization', element: <HubspotAuthorizations/> },
  { path: '/hubspot-authorizations', name: 'Hubspot Authorizations', element: <HubspotAuthorizations/> },
  { path: '/salesforce-authorizations/:itemUuid', name: 'Salesforce Authorizations', element: <SalesforceAuthorizations/> },
  { path: '/salesforce-authorizations', name: 'Salesforce Authorizations', element: <SalesforceAuthorizations/> },
  { path: '/settings', name: 'Settings', element: <Settings/> },
  { path: '/env-variables/:itemUuid', name: 'Environment variables', element: <EnvVariables/> },
  { path: '/env-variables', name: 'Environment variables', element: <EnvVariables/> },
  { path: '/metrics-old/:itemUuid', name: 'Metrics', element: <Metrics/> },
  { path: '/metrics-old', name: 'Metrics', element: <Metrics/> },
  { path: '/api-clients', name: 'API Clients', element: <APIClients/> },
  { path: '/users', name: 'Users', element: <Users/> },
  { path: '/slow-tasks/:dateFrom/:dateTo/:threshold/:taskType', name: 'Slow tasks', element: <SlowTasks/> },
  { path: '/slow-tasks/:dateFrom/:dateTo/:threshold', name: 'Slow tasks', element: <SlowTasks/> },
  { path: '/slow-tasks/:dateFrom/:dateTo', name: 'Slow tasks', element: <SlowTasks/> },
  { path: '/slow-tasks', name: 'Slow tasks', element: <SlowTasks/> },
  //{ name: '404', element: Page404 },
];

export default routes;
