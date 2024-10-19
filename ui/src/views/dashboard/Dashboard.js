import React from 'react';
import Tour from 'reactour';
import {connect} from "react-redux";
import { Navigate } from 'react-router-dom';
import { DateTime } from "luxon";

import Button from 'react-bootstrap/Button';
import Dropdown from "react-bootstrap/Dropdown";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Card from 'react-bootstrap/Card';
import Spinner from "react-bootstrap/Spinner";

import {tourSteps} from "../../actions/tourSteps";
import {releaseNotes} from "../../actions/releaseNotes";

import {PERMISSION_LOGGED_IN, validatePermissions} from "../../actions/pipeline";

import PipelineWidget from "../../components/widgets/pipelineWidget.component";
import IntegrationWidget from "../../components/widgets/integrationWidget.component";

import PolarAreaChart, {polarAreaDashboardData} from "../../components/charts/polarAreaChart.component";
import LineChart, {lineChartOptions} from "../../components/charts/lineChart.component";
import BarChart, {barDataSample} from "../../components/charts/barChart.component";

import PipelineService, {PIPELINE_TYPE_METRIC, PIPELINE_TYPE_GENERIC} from "../../services/pipeline.service";

const ReleaseVersion = 17;

const isDarkThemeEnabled = localStorage.getItem('darkTheme') !== "false";

const stateLineChartData = {
  labels: [],
  datasets: [
    {
      label: '',
      data: [],
      borderColor: isDarkThemeEnabled ? 'rgb(206, 150, 250)' : 'rgb(176, 96, 239)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

const timeFrames = [ 
    "week", 
    "month",
];

class Dashboard extends React.Component {
    constructor(props) {
        super(props);

        this.handleGetPipelines = this.handleGetPipelines.bind(this);
        this.handleGetDashboard = this.handleGetDashboard.bind(this);
        this.onChangePipelineGroupBy = this.onChangePipelineGroupBy.bind(this);
        this.onChangeMetricGroupBy = this.onChangeMetricGroupBy.bind(this);

        let lChartOptions = JSON.parse(JSON.stringify(lineChartOptions));

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isTourOpen:
                window.innerWidth > 992
                && !localStorage.getItem('tourIsWatched'),
            isReleaseNotesOpen:
                window.innerWidth > 992
                && localStorage.getItem('tourIsWatched')
                && (!localStorage.getItem('releaseNotesWatched') || localStorage.getItem('releaseNotesWatched') < ReleaseVersion),
            transparent: false,
            pipelines: null,
            dashboard: null,
            pipelineChartData: JSON.parse(JSON.stringify(polarAreaDashboardData)),
            metricChartData: JSON.parse(JSON.stringify(polarAreaDashboardData)),
            newGroupedPipelines: null,
            newGroupedMetrics: null,
            pipelineRunsPerMonth: null,
            metricRunsPerMonth: null,
            pipelineLineChartData: JSON.parse(JSON.stringify(stateLineChartData)),
            metricLineChartData: JSON.parse(JSON.stringify(stateLineChartData)),
            pipelineRunsChartData: JSON.parse(JSON.stringify(stateLineChartData)),
            metricRunsChartData: JSON.parse(JSON.stringify(stateLineChartData)),
            pipelineGroupBy: "month",
            metricGroupBy: "month",
            lChartOptions: lChartOptions,
        };
    }

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    onChangePipelineGroupBy = async(period) => {
        await this.promisedSetState({
            pipelineGroupBy: period,
            newGroupedPipelines: null,
            pipelineLineChartData: JSON.parse(JSON.stringify(stateLineChartData)),
        });

        await this.handleGetGroupedItems(
            this.state.organization.uuid,
            PIPELINE_TYPE_GENERIC,
            this.state.pipelineGroupBy 
        );
    };

    onChangeMetricGroupBy = async(period) => {
        await this.promisedSetState({
            metricGroupBy: period,
            newGroupedMetrics: null,
            metricLineChartData: JSON.parse(JSON.stringify(stateLineChartData)),
        });

        await this.handleGetGroupedItems(
            this.state.organization.uuid,
            PIPELINE_TYPE_METRIC,
            this.state.metricGroupBy 
        );
    };

    closeTour = () => {
        localStorage.setItem('tourIsWatched', true);
        localStorage.setItem('releaseNotesWatched', ReleaseVersion);
        window.location.reload();
    };

    closeReleaseNotes = () => {
        this.setState({isReleaseNotesOpen: false});
        localStorage.setItem('releaseNotesWatched', ReleaseVersion);
    };

    componentDidMount() {
        document.title = 'Ylem'
        this.handleGetPipelines(this.state.organization.uuid);
        this.handleGetDashboard(this.state.organization.uuid);

        this.handleGetGroupedItems(
            this.state.organization.uuid,
            PIPELINE_TYPE_GENERIC,
            this.state.pipelineGroupBy 
        );

        this.handleGetGroupedItems(
            this.state.organization.uuid,
            PIPELINE_TYPE_METRIC,
            this.state.metricGroupBy 
        );

        this.handleGetRunsPerMonth(
            this.state.organization.uuid,
            PIPELINE_TYPE_GENERIC
        );

        this.handleGetRunsPerMonth(
            this.state.organization.uuid,
            PIPELINE_TYPE_METRIC
        );
    };

    handleGetRunsPerMonth = async(uuid, type) => {
        let items = (type === PIPELINE_TYPE_GENERIC) 
            ? this.state.pipelineRunsPerMonth
            : this.state.metricRunsPerMonth;

        let chartData = (type === PIPELINE_TYPE_GENERIC) 
            ? this.state.pipelineRunsChartData 
            : this.state.metricRunsChartData;

        chartData = JSON.parse(JSON.stringify(barDataSample));

        if (
            items === null
            || items.length === 0
        ) {
            items = PipelineService.getRunsPerMonthByOrganization(uuid, type);

            await Promise.resolve(items)
                .then(async(items) => {
                    if (items.data && items.data.items && items.data.items !== null) {
                        items.data.items.sort().reverse()
                        
                        for(var i = 0; i < items.data.items.length; i++){
                            chartData.datasets[0].data.push(items.data.items[i].run_count);
                            chartData.labels.push(
                                DateTime.fromFormat(items.data.items[i].year_month, 'yyyy-MM').toLocaleString({ month: 'short', year: 'numeric' })
                            );
                        }

                        chartData.datasets.splice(1, 1);

                        if (type === PIPELINE_TYPE_GENERIC) {
                            chartData.datasets[0].label = " Pipeline runs"
                            await this.promisedSetState({
                                pipelineRunsPerMonth: items.data.items,
                                pipelineRunsChartData: chartData,
                            });
                        } else {
                            chartData.datasets[0].label = " Metric runs"
                            await this.promisedSetState({
                                metricRunsPerMonth: items.data.items,
                                metricRunsChartData: chartData,
                            });
                        }
                    } else {
                        if (type === PIPELINE_TYPE_GENERIC) {
                            await this.promisedSetState({pipelineRunsPerMonth: []});
                        } else {
                            await this.promisedSetState({metricRunsPerMonth: []});
                        }
                    }
                })
                .catch(async() => {
                    if (type === PIPELINE_TYPE_GENERIC) {
                        await this.promisedSetState({pipelineRunsPerMonth: []});
                    } else {
                        await this.promisedSetState({metricRunsPerMonth: []});
                    }
                });
        }
    };

    handleGetGroupedItems = async(uuid, type, groupBy) => {
        let items = (type === PIPELINE_TYPE_GENERIC) 
            ? this.state.newGroupedPipelines 
            : this.state.newGroupedMetrics;

        let chartData = (type === PIPELINE_TYPE_GENERIC) 
            ? this.state.pipelineLineChartData 
            : this.state.metricLineChartData;

        if (
            items === null
            || items.length === 0
        ) {
            items = PipelineService.getNewGroupedItemsByOrganization(uuid, type, groupBy);

            await Promise.resolve(items)
                .then(async(items) => {
                    if (items.data && items.data.items && items.data.items !== null) {
                        for(var i = 0; i < items.data.items.length; i++){
                            chartData.datasets[0].data.push(items.data.items[i].count);
                            chartData.labels.push(
                                items.data.items[i].week === 0
                                    ? items.data.items[i].month + " " + items.data.items[i].year
                                    : 'W' + items.data.items[i].week + " " + items.data.items[i].month + " " + items.data.items[i].year
                            );
                        }

                        if (type === PIPELINE_TYPE_GENERIC) {
                            chartData.datasets[0].label = " New pipelines"
                            await this.promisedSetState({
                                newGroupedPipelines: items.data.items,
                                pipelineLineChartData: chartData,
                            });
                        } else {
                            chartData.datasets[0].label = " New metrics"
                            await this.promisedSetState({
                                newGroupedMetrics: items.data.items,
                                metricLineChartData: chartData,
                            });
                        }
                    } else {
                        if (type === PIPELINE_TYPE_GENERIC) {
                            await this.promisedSetState({newGroupedPipelines: []});
                        } else {
                            await this.promisedSetState({newGroupedMetrics: []});
                        }
                    }
                })
                .catch(async() => {
                    if (type === PIPELINE_TYPE_GENERIC) {
                        await this.promisedSetState({newGroupedPipelines: []});
                    } else {
                        await this.promisedSetState({newGroupedMetrics: []});
                    }
                });
        }
    };

    handleGetPipelines = async(uuid) => {
        let pipelines = this.state.pipelines;

        if (
            pipelines === null
            || pipelines.length === 0
        ) {
            pipelines = PipelineService.getPipelinesByOrganization(uuid);

            await Promise.resolve(pipelines)
                .then(async(pipelines) => {
                    if (pipelines.data && pipelines.data.items && pipelines.data.items !== null) {
                        var items = pipelines.data.items.filter(k => k.type === PIPELINE_TYPE_GENERIC);
                        
                        await this.promisedSetState({pipelines: items});
                        if (items.length > 0) {
                            await this.promisedSetState({activePipeline: items[0]});
                        }
                    } else {
                        await this.promisedSetState({pipelines: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({pipelines: []});
                });
        }
    };

    handleGetDashboard = async(uuid) => {
        let dashboard = this.state.dashboard;

        if (dashboard === null) {
            dashboard = PipelineService.getDashboardByOrganization(uuid);

            await Promise.resolve(dashboard)
                .then(async(dashboard) => {
                    if (dashboard && dashboard.data !== null) {
                        let pipelineChartData = this.state.pipelineChartData;
                        pipelineChartData.labels = [
                            ' Active pipelines',
                            ' My pipelines',
                            ' Scheduled pipelines',
                            ' Externally triggered pipelines',
                            ' New pipelines', 
                            ' Recently updated pipelines',
                            ' Pipeline templates',
                        ];
                        pipelineChartData.datasets[0].data = [
                            dashboard.data.num_active_pipelines,
                            dashboard.data.num_my_pipelines,
                            dashboard.data.num_scheduled_pipelines,
                            dashboard.data.num_externally_triggered_pipelines,
                            dashboard.data.num_new_pipelines,
                            dashboard.data.num_recently_updated_pipelines,
                            dashboard.data.num_pipeline_templates,
                        ];

                        let metricChartData = this.state.metricChartData;
                        metricChartData.labels = [
                            ' Active metrics', 
                            ' New metrics',
                            ' My metrics', 
                            ' Recently updated metrics',
                            ' Metric templates',
                        ];
                        metricChartData.datasets[0].data = [
                            dashboard.data.num_active_metrics,
                            dashboard.data.num_new_metrics,
                            dashboard.data.num_my_metrics,
                            dashboard.data.num_recently_updated_metrics,
                            dashboard.data.num_metric_templates,
                        ];

                        await this.promisedSetState({
                            dashboard: dashboard.data,
                            pipelineChartData: pipelineChartData,
                            metricChartData: metricChartData,
                        });
                    } else {
                        await this.promisedSetState({dashboard: {}});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({dashboard: {}});
                });
        }
    };

    render() {
        const { isLoggedIn, user } = this.props;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
            return <Navigate to="/login" />;
        }

    	const { 
            isTourOpen, 
            isReleaseNotesOpen, 
            dashboard, 
            pipelineChartData, 
            metricChartData,
            pipelineLineChartData,
            metricLineChartData,
            pipelineRunsChartData,
            metricRunsChartData,
            newGroupedPipelines,
            newGroupedMetrics,
            pipelineRunsPerMonth,
            metricRunsPerMonth,
            pipelineGroupBy,
            metricGroupBy,
            lChartOptions,
        } = this.state;

  		return (
  			<>
                { dashboard !== null
                    ?
                        <div className="dashboard">
                            <h2 className="tasksDropdownPreTitle">Pipelines</h2><br/>
                            <Row className="mb-5">
                                <Col sm={6} className="px-5 dashboardChart">
                                    <PolarAreaChart data={pipelineChartData} />
                                </Col>
                                <Col sm={3} className="pt-4 px-5">
                                    <span className="statValue">
                                        {dashboard.num_active_pipelines}
                                    </span><br/>
                                    <span className="statKey">Active pipelines</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_scheduled_pipelines}
                                    </span><br/>
                                    <span className="statKey">Scheduled pipelines</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_externally_triggered_pipelines}
                                    </span><br/>
                                    <span className="statKey">Externally triggered pipelines</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_pipeline_templates}
                                    </span><br/>
                                    <span className="statKey">Pipeline templates</span><br/><br/>
                                </Col>
                                <Col sm={3} className="pt-4 px-5">
                                    <span className="statValue">
                                        {dashboard.num_my_pipelines}
                                    </span><br/>
                                    <span className="statKey">My pipelines</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_new_pipelines}
                                    </span><br/>
                                    <span className="statKey">New pipelines</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_recently_updated_pipelines}
                                    </span><br/>
                                    <span className="statKey">Recently updated pipelines</span><br/><br/>
                                </Col>
                            </Row>

                            { pipelineRunsPerMonth !== null
                                ?
                                    <div className="mb-5">
                                        <h2 className="tasksDropdownPreTitle px-3">Pipeline runs per month</h2>
                                        <div className="dashboardLineChart w-98-left">
                                            <BarChart
                                                data={pipelineRunsChartData}
                                            />
                                        </div>
                                    </div>
                                : <div className="text-center dashboardLineChart"><Spinner animation="grow" className="spinner-primary"/></div>
                            }

                            { newGroupedPipelines !== null
                                ?
                                    <div className="mb-5">
                                        <h2 className="tasksDropdownPreTitle px-3 float-left">New pipelines by </h2>
                                        <Dropdown size="lg" className="mb-4 mt-2 float-left">
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className=""
                                            >
                                                {pipelineGroupBy.charAt(0).toUpperCase() + pipelineGroupBy.slice(1)}
                                            </Dropdown.Toggle>

                                            <Dropdown.Menu className="tasks">
                                                {timeFrames.map(value => (
                                                    <Dropdown.Item
                                                        value={value}
                                                        key={value}
                                                        active={value === pipelineGroupBy}
                                                        onClick={(e) => this.onChangePipelineGroupBy(value)}
                                                    >
                                                        {value.charAt(0).toUpperCase() + value.slice(1)}
                                                    </Dropdown.Item>
                                                ))}
                                            </Dropdown.Menu>
                                        </Dropdown>
                                        <div className="clearfix"></div>
                                        <div className="dashboardLineChart w-98-left">
                                            <LineChart
                                                data={pipelineLineChartData}
                                                options={lChartOptions}
                                            />
                                        </div>
                                    </div>
                                : <div className="text-center dashboardLineChart"><Spinner animation="grow" className="spinner-primary"/></div>
                            }

                            <h2 className="tasksDropdownPreTitle">Metrics</h2><br/>
                            <Row className="mb-5">
                                <Col sm={6} className="px-5 dashboardChart">
                                    <PolarAreaChart data={metricChartData} />
                                </Col>
                                <Col sm={3} className="pt-5 px-5">
                                    <span className="statValue">
                                        {dashboard.num_active_metrics}
                                    </span><br/>
                                    <span className="statKey">Active metrics</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_my_metrics}
                                    </span><br/>
                                    <span className="statKey">My metrics</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_metric_templates}
                                    </span><br/>
                                    <span className="statKey">Metric templates</span>
                                </Col>
                                <Col sm={3} className="pt-5 px-5">
                                    <span className="statValue">
                                        {dashboard.num_new_metrics}
                                    </span><br/>
                                    <span className="statKey">New metrics</span><br/><br/>
                                    <span className="statValue">
                                        {dashboard.num_recently_updated_metrics}
                                    </span><br/>
                                    <span className="statKey">Recently updated metrics</span><br/><br/>
                                </Col>
                            </Row>

                            { metricRunsPerMonth !== null
                                ?
                                    <div className="mb-5">
                                        <h2 className="tasksDropdownPreTitle px-3">Metric runs per month</h2>
                                        <div className="dashboardLineChart w-98-left">
                                            <BarChart
                                                data={metricRunsChartData}
                                            />
                                        </div>
                                    </div>
                                : <div className="text-center dashboardLineChart"><Spinner animation="grow" className="spinner-primary"/></div>
                            }

                            { newGroupedMetrics !== null
                                ?
                                    <div className="mb-5">
                                        <h2 className="tasksDropdownPreTitle px-3 float-left">New metrics by </h2>
                                        <Dropdown size="lg" className="mb-4 mt-2 float-left">
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className=""
                                            >
                                                {metricGroupBy.charAt(0).toUpperCase() + metricGroupBy.slice(1)}
                                            </Dropdown.Toggle>

                                            <Dropdown.Menu className="tasks">
                                                {timeFrames.map(value => (
                                                    <Dropdown.Item
                                                        value={value}
                                                        key={value}
                                                        active={value === metricGroupBy}
                                                        onClick={(e) => this.onChangeMetricGroupBy(value)}
                                                    >
                                                        {value.charAt(0).toUpperCase() + value.slice(1)}
                                                    </Dropdown.Item>
                                                ))}
                                            </Dropdown.Menu>
                                        </Dropdown>
                                        <div className="clearfix"></div>
                                        <div className="dashboardLineChart w-98-left"> 
                                            <LineChart
                                                data={metricLineChartData}
                                                options={lChartOptions}
                                            />
                                        </div>
                                    </div>
                                : <div className="text-center dashboardLineChart"><Spinner animation="grow" className="spinner-primary"/></div>
                            }
                        </div>
                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }

                <h2 className="tasksDropdownPreTitle">Latest updates</h2><br/>

                <Row className="mb-4">
                    <Col sm={4}>
                        <Card className="withHeader dashboardCard">
                            <Card.Header>
                                Last pipelines
                            </Card.Header>
                            <Card.Body>
                                <PipelineWidget type={PIPELINE_TYPE_GENERIC}/>
                            </Card.Body>
                        </Card>
                    </Col>
                    <Col sm={4}>
                        <Card className="withHeader dashboardCard">
                            <Card.Header>
                                Last Metrics
                            </Card.Header>
                            <Card.Body>
                                <PipelineWidget type={PIPELINE_TYPE_METRIC}/>
                            </Card.Body>
                        </Card>
                    </Col>
                    <Col sm={4}>
                        <Card className="withHeader dashboardCard">
                            <Card.Header>
                                Integrations
                            </Card.Header>
                            <Card.Body>
                                <IntegrationWidget/>
                            </Card.Body>
                        </Card>
                    </Col>
                </Row>

                 <Tour
                    steps={tourSteps}
                    isOpen={isTourOpen}
                    showNavigationNumber={false}
                    showNumber={false}
                    highlightedMaskClassName="tourMask"
                    rounded={3}
                    maskSpace={0}
                    onRequestClose={this.closeTour} 
                />
                <Tour
                    steps={releaseNotes}
                    isOpen={isReleaseNotesOpen}
                    showNavigationNumber={false}
                    showNumber={false}
                    highlightedMaskClassName="tourMask"
                    rounded={3}
                    maskSpace={0}
                    lastStepNextButton={<Button color="primary">Let's give it a try!</Button>}
                    onRequestClose={this.closeReleaseNotes} 
                />
            </>
  		)
  	}
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    return {
        isLoggedIn,
        user,
    };
}

export default connect(mapStateToProps)(Dashboard);
