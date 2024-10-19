import React from 'react'
import { Link } from 'react-router-dom';

export const SalesforceAuthorizationInfo = {
    title: "Salesforce Authorizations",
    content: 
        <>
            <div>
                To be able to create cases in a real-time streaming mode in <strong>Salesforce</strong> of your organization, you need to authorize Ylem to do so first.
                <br/><br/>
                It may happen that you don't have admin permissions in your Salesforce. In this case, you need to check it with the <strong>administrator</strong> of your organization.
                <br/><br/>
                Please note, that a user must be assigned to the system permissions "API Enabled", "Apex REST Services".
            </div>
        </>
}

export const SlowTasksInfo = {
    title: "Slow Tasks",
    content: 
        <>
            <div>
                Here on the profiling page, you can filter the historical log of tasks by their type, execution date, and duration threshold.<br/><br/> 
                It helps you to identify slow tasks (especially database queries) and start working on optimization.
            </div>
        </>
}

export const SlackAuthorizationInfo = {
    title: "Slack Authorizations",
    content:
        <>
            <div>
                To be able to stream messages in <strong>Slack</strong> channels of your organization, you need to authorize Ylem to have access to it first.
                <br/><br/>
                It may happen that you don't have admin permissions in your Slack. In this case, you need to check it with the <strong>administrator</strong> of your organization.
            </div>
        </>
}

export const EnvVariablesInfo = {
    title: "Environment Variables",
    content:
        <>
            <div>
                Environment variables allow you to set variables in one place and reuse it in pipelines and metrics as <strong>ENV_variable_name</strong>. 
                <br/><br/>
                For example, if you have a certain KPI threshold equal to <strong>100</strong>, you can set its variable as <strong>key="KPI_THRESHOLD"</strong> and <strong>value=100</strong> here and reuse it everywhere else as <strong>ENV_KPI_THRESHOLD</strong>.
                <br/><br/>
                Then if you want to make it <strong>110</strong> instead of <strong>100</strong>, you can do it in one place, instead of changing all of your pipelines and metrics.
            </div>
        </>
}

export const MetricsInfo = {
    title: "Metrics",
    content:
        <>
            <div>
                The main purpose of metrics is to give you a simple interface to monitor your KPIs, SLAs, and other metrics and run various pipelines based on their value.<br/><br/>
                It may also happen that you already have metrics calculated on your data warehouse level. In this case, you can just retrieve them here with the SQL query, create a monitoring schedule, and thresholds and execute Ylem's pipeline based on their value.
            </div>
        </>
}

export const JiraAuthorizationInfo = {
    title: "Jira Authorizations",
    content:
        <>
            <div>
                To be able to create issues in your <strong>Jira Cloud</strong> in real-time, you need to authorize Ylem to do so first.<br/><br/> 
                Please note, that Jira Server at the moment is not supported and we can only stream to Jira Cloud.
            </div>
        </>
};

export const HubspotAuthorizationInfo = {
    title: "Hubspot Authorizations",
    content:
        <>
            <div>
                To be able to create tickets in your <strong>Hubspot</strong> in real-time, you need to authorize Ylem to do so first.<br/><br/>
                It may happen that you don't have admin permissions in your Hubspot. In this case, you need to check it with the <strong>administrator</strong> of your organization.
            </div>
        </>
};

export const PipelinesInfo = {
    title: "Pipelines",
    content:
        <>
            <div>
                <strong>Pipeline</strong> is the heart of Ylem. It orchestrates the entire streaming process from triggering to sending data to a destination.
                <br/><br/>
                With various types of tasks, you can retrieve data, enrich it, merge, transform, ingest, convert to various formats, and much more.
                <br/><br/>
                <strong>Queries</strong> are usually the first ones in the pipeline. Here you can write an SQL query you want to run against your <Link to="/data-sources">Data Sources</Link> and a cronjob-like schedule when to run it. The result of the query can be sent to Aggregators, Conditions, or Transformers. 
                <br/><br/>
                <strong>Aggregators</strong> help you to implement basic aggregating formulas and pre-build functions such as COUNT(), AVG(), SUM(), and other ones.
                <br/><br/>
                <strong>Condition</strong> is a block where you can compare input and expectations and continue with different pipeline scenarios depending on if it is true or false.
                <br/><br/>
                <strong>Merge</strong> allows merging data coming from two and more data sources.
                <br/><br/>
                <strong>For Each</strong> allows running the next task for each of the items separately.
                <br/><br/>
                <strong>Transformers</strong> help you to transform data in different formats or extract some values.
                <br/><br/>
                <strong>External trigger</strong> is a task that you place first into your workkflow if you want to trigger it from a data streaming platform like Kafka.
                <br/><br/>
                <strong>Processor</strong> is a task that allows you to transform, filter, and map data using the functionality of JQ library.
                <br/><br/>
                In the end, you can either call APIs with <strong>API Calls</strong> or send notifications with <strong>Notification</strong> tasks to the <Link to="/integrations">Integrations</Link> you've configured before.
            </div>
        </>
};

export const DataSourcesInfo = {
    title: "Data Sources",
    content:
        <>
            <div>
                Here you can configure your data sources to be used in pipelines. It can be your production database replica, your data warehouse, staging or testing database, or any other storage.
                <br/><br/>
                We recommend having both production and testing connections ready, so you can always properly test your pipelines before using them on production.
                <br/><br/>
                At the moment we support a wide range of <strong>SQL-based</strong> data sources, such as MySQL, Snowflake, Postgres, Googgle Big Query, AWS RDS, and others. And this list is growing since we are adding new ones.
            </div>
        </>
};

export const IntegrationsInfo = {
    title: "Integrations",
    content:
        <>
            <div>
                Integrations are the external channels from where you can read data for your streaming and where you can also navigate your data streams.
                <br/><br/>
                Integrations are divided into the following categories:
                <br/><br/>
                <strong>Read-write</strong><br/>
                <ul>
                    <li>APIs</li>
                    <li>SQL-databases (MySQL, PostgreSQL, Snowflake, AWS RDS, Google Big Query, etc.)</li>
                </ul>

                <strong>Write</strong><br/>
                <ul>
                    <li>Slack</li>
                    <li>Atlassian Jira</li>
                    <li>Incident.io</li>
                    <li>Opsgenie</li>
                    <li>Google Sheets</li>
                    <li>etc.</li>
                </ul>

                <strong>Read</strong><br/>
                <ul>
                    <li>Elasticsearch</li>
                </ul>
                On top of that, we actively support streaming from and to Apache Kafka, RabbitMQ, Google Pub/Sub, and other similar technologies that you don't need to configure here but can use through our API and open-source plugins from our repository.
                <br/><br/>
                More information can be found in our documentation: <a href="https://docs.ylem.co/datamin-api/api-endpoints" target="_blank" rel="noreferrer">docs.ylem.co</a>.
            </div>
        </>
};

export const SettingsInfo = {
    title: "Settings",
    content:
        <>
            <div>
                The <strong>settings</strong> page allows you to configure your personal settings including your passwords, names, and email address.
                <br/><br/>
                If you have administrative permissions, here you can also set up your <strong>organization</strong>.
            </div>
        </>
};

export const UsersInfo = {
    title: "Users",
    content:
        <>
            <div>
                Here you can invite other team members to join your organization at Ylem. Just add their Emails to the list and we send them invitation links.
                <br/><br/>
                As soon as they accept the invitation, you can make the <strong>administrators</strong> or keep them as <strong>team members</strong>. Administrators can invite other users.
            </div>
        </>
};

export const OAuthInfo = {
    title: "API OAuth Clients",
    content:
        <>
            <div>
                <strong>API OAuth Clients</strong> are needed in case you want to call execution of pipelines from outside of Ylem via API.
                <br/><br/>
                To do that, you need to create a client first and <strong>copy and save</strong> its client secret. You won't be able to see it again in the UI.
                <br/><br/>
                The available endpoints are and how to use them is described in our documentation: <a href="https://docs.ylem.co/datamin-api/api-endpoints" target="_blank" rel="noreferrer">docs.ylem.co</a>
            </div>
        </>
};
