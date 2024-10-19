export const tourSteps = [
    {
        content: () => (
            <div className="text-center">
                Welcome to <strong>Ylem - the data streaming platform!</strong><br/><br/>Let me show you around. It won't take more than 1 minute of your time.<br/>
                <div className="text-center">
                    <img alt="Welcome to Ylem" src="/images/tour/welcome.png" width="400px"/>
                </div>
            </div>
        ),
        className: 'tour',
    },
    {
        selector: '.tour-step-dashboard',
        content: () => (
            <div>
                Our main goal is to make data streaming easy, reliable, affordable, and as close to real-time as possible.
                <br/><br/>
                To achieve that for your organization with Ylem, you will configure <strong>integrations</strong> with your APIs, data sources, and other services and orchestrate <strong>pipelines</strong> to move data between them.
                <br/><br/>
                <div className="text-center">
                    <img alt="Welcome to Ylem" src="/images/tour/dashboard.png" width="600px"/>
                </div>
            </div>
        ),
        maskClassName: 'tourMask',
    },
    {
        selector: '.tour-step-integrations',
        content: () => (
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
        ),
        position: [50, 20],
    },
    {
        selector: '.tour-step-pipelines',
        content: () => (
            <div>
                Pipeline is the heart of Ylem. It orchestrates the entire streaming process from triggering to sending data to a destination.
                <br/><br/>
                With various types of tasks, you can retrieve data, enrich it, merge, transform, ingest, convert to various formats, and much more. Give it a try yourself.
                <div className="text-center">
                    <img alt="Pipelines" src="/images/tour/pipeline.png" width="700px" className="tourImage"/>
                </div>
            </div>
        ),
        position: [50, 30],
    },
    {
        selector: '.tour-step-metrics',
        content: () => (
            <div>
                The main purpose of metrics is to give you a simple interface to monitor your KPIs, SLAs, and other metrics and run various pipelines based on their value.<br/><br/>
                It may also happen that you already have metrics calculated on your data warehouse level. In this case, you can just retrieve them here with the SQL query, create a monitoring schedule, and thresholds and execute Ylem's pipeline based on their value.
                <div className="text-center">
                    <img alt="Pipelines" src="/images/tour/metrics.png" width="700px" className="tourImage"/>
                </div>
            </div>
        ),
        position: [50, 30],
    },
    {
        selector: '.tour-step-oauth-clients',
        content: () => (
            <div>
                Ylem is an <strong>API-driven</strong> platform that allows you to trigger your pipelines in multiple flexible ways:
                <ul>
                    <li>By schedule,</li>
                    <li>Manually from the UI,</li>
                    <li>From data streaming solutions (<strong>Kafka, RabbitMQ, Google Pub/Sub, Amazon SQS</strong>) or serverless solutions (<strong>AWS Lambda, Google Cloud Functions</strong>) using our open-source libraries,</li> 
                    <li>Using our powerful API.</li>
                </ul>
                For the last three options you will need OAuth clients which you can configure here:
                <div className="text-center">
                    <img alt="OAuth" src="/images/tour/oauth.png" width="500px" className="tourImage"/>
                </div>
            </div>
        ),
        position: [50, 30],
    },
    {
        selector: '.tour-step-integrations',
        content: () => (
            <div>
                Environment variables allow you to set variables in one place and reuse them in pipelines and metrics as <strong>ENV_variable_name</strong>.
                <div className="text-center">
                    <img alt="Env variables" src="/images/tour/env-variables.png" width="700px" className="tourImage"/>
                </div>
            </div>
        ),
        position: [50, 30],
    },
    {
        selector: '.tour-step-profiling',
        content: () => (
            <div>
                Last but not least, with our super-powered <strong>profiling</strong> and <strong>logging</strong> system, you will always be aware of slow tasks and queries, failed pipeline runs, and incorrect inputs/outputs.
                <div className="text-center">
                    <img alt="Profiling" src="/images/tour/profiling.png" width="700px" className="tourImage"/>
                </div>
            </div>
        ),
        position: [50, 30],
    },
    {
        content: () => (
            <div className="text-center">
                Thanks for watching! And have a successful data streaming orchestration journey with Ylem!<br/><br/>
                <div className="text-center">
                    <img alt="Welcome to Ylem" src="/images/tour/welcome.png" width="400px"/>
                </div>
            </div>
        ),
    },
]
