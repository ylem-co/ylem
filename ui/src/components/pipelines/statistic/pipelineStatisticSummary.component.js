import React from 'react';
import { DateTime } from "luxon";
import prettyMilliseconds from 'pretty-ms';

import StatService from "../../../services/stat.service";
import { PIPELINE_TYPE_METRIC } from "../../../services/pipeline.service";


import Spinner from "react-bootstrap/Spinner";
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Dropdown from "react-bootstrap/Dropdown";

import BarChart, {barDataSample} from "../../charts/barChart.component";
import LineChart, {lineDataSample, lineValueDataSample, lineChartOptions, lineChartValueOptions} from "../../charts/lineChart.component";

const isDarkThemeEnabled = localStorage.getItem('darkTheme') !== "false";

const timeFrames = [
    "day", 
    "week", 
    "month", 
    "quarter", 
    "year"
];

const emptyStateStats = {
    "num_of_successes": 0,
    "num_of_failures": 0,
    "average_duration": 0,
    "is_last_run_successful": null,
    "last_run_duration": 0,
    "last_run_executed_at": DateTime.now().toSQL(),
};

const emptyStateValues = [{
    "pipeline_uuid": "63e69014-a318-42f9-b741-5614bd9a4181",
    "metric_value": 20,
    "executed_at": "2023-01-01 00:00:00"
}, {
    "pipeline_uuid": "63e69014-a318-42f9-b741-5614bd9a4181",
    "metric_value": 21,
    "executed_at": "2023-02-01 00:00:00"
}, {
    "pipeline_uuid": "63e69014-a318-42f9-b741-5614bd9a4181",
    "metric_value": 15,
    "executed_at": "2023-03-01 00:00:00"
}, {
    "pipeline_uuid": "63e69014-a318-42f9-b741-5614bd9a4181",
    "metric_value": 30,
    "executed_at": "2023-04-01 00:00:00"
}];

const emptyStateLineChartData = {
  labels: [
    "Jan 2023",
    "Feb 2023",
    "Mar 2023",
    "Apr 2023",
    "May 2023",
    "Jun 2023",
    "Jul 2023",
    "Aug 2023",
  ],
  datasets: [
    {
      label: ' Average duration',
      data: [300, 2000, 345, 570, 190, 59, 890, 60],
      backgroundColor: 'rgb(128, 128, 128, .3)',
      color: 'rgb(128, 128, 128, .3)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

const emptyStateLineChartValueData = {
  labels: [
    "Jan 2023",
    "Feb 2023",
    "Mar 2023",
    "Apr 2023",
    "May 2023",
    "Jun 2023",
    "Jul 2023",
    "Aug 2023",
  ],
  datasets: [
    {
      label: ' Value',
      data: [23, 39, 17, 19, 10, 32, 33, 37],
      backgroundColor: 'rgb(128, 128, 128, .3)',
      color: 'rgb(128, 128, 128, .3)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

const emptyStateBarChartData = {
  labels: [
    "Jan 2023",
    "Feb 2023",
    "Mar 2023",
    "Apr 2023",
    "May 2023",
    "Jun 2023",
    "Jul 2023",
    "Aug 2023",
  ],
  datasets: [
    {
      label: 'Successful runs',
      data: [10, 20, 30, 60, 10, 70, 20, 60],
      backgroundColor: 'rgb(128, 128, 128, .3)',
    },
    {
      label: 'Failures',
      data: [3, 2, 3, 10, 10, 20, 5, 20],
      backgroundColor: 'rgb(192, 192, 192, .4)',
    },
  ],
};

const emptyStateLineChartDataZeros = {
  labels: [
    "Jan 2023",
    "Feb 2023",
    "Mar 2023",
    "Apr 2023",
    "May 2023",
    "Jun 2023",
    "Jul 2023",
    "Aug 2023",
  ],
  datasets: [
    {
      label: ' Average duration',
      data: [0, 0, 0, 0, 0, 0, 0, 0],
      backgroundColor: 'rgb(128, 128, 128, .2)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

const emptyStateLineChartValueDataZeros = {
  labels: [
    "Jan 2023",
    "Feb 2023",
    "Mar 2023",
    "Apr 2023",
    "May 2023",
    "Jun 2023",
    "Jul 2023",
    "Aug 2023",
  ],
  datasets: [
    {
      label: ' Value',
      data: [0, 0, 0, 0, 0, 0, 0, 0],
      backgroundColor: 'rgb(128, 128, 128, .2)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

const emptyStateBarChartDataZeros = {
  labels: [
    "Jan 2023",
    "Feb 2023",
    "Mar 2023",
    "Apr 2023",
    "May 2023",
    "Jun 2023",
    "Jul 2023",
    "Aug 2023",
  ],
  datasets: [
    {
      label: 'Successful runs',
      data: [0, 0, 0, 0, 0, 0, 0, 0],
      backgroundColor: 'rgb(128, 128, 128, .3)',
      color: isDarkThemeEnabled ? 'white' : 'black',
    },
    {
      label: 'Failures',
      data: [0, 0, 0, 0, 0, 0, 0, 0],
      backgroundColor: 'rgb(192, 192, 192, .4)',
      color: isDarkThemeEnabled ? 'white' : 'black',
    },
  ],
};

const emptyStateAggregatedStats = [{
    "empty": "state"
}];

class PipelineStatisticSummary extends React.Component {
    constructor(props) {
        super(props);

        let dF = this.props.dateFrom;
        if (dF) {
            let dFParts = dF.split(' ');
            if (dFParts.length === 1) {
                dF = dF + " 00:00:00";
            }
        }

        let dT = this.props.dateTo;
        if (dT) {
            let dTParts = dT.split(' ');
            if (dTParts.length === 1) {
                dT = dT + " 23:59:59";
            }
        }

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            item: this.props.item,
            type: this.props.type,
            dateFrom: dF ||  DateTime.now().plus({ days: -7 }).toSQLDate() + " 00:00:00",
            dateTo: dT || DateTime.now().toSQLDate() + " 23:59:59",
            stats: null,
            values: null,
            barChartData: null,
            lineChartData: null,
            lineChartValueData: null,
            aggregatedStats: null,
            period: 'day',
            periodLength: 0,
            lChartOptions: JSON.parse(JSON.stringify(lineChartOptions)),
            lChartValueOptions: JSON.parse(JSON.stringify(lineChartValueOptions)),
        };
    }

    getAggregatedDateFrom = (dateTo, period, periodLength) => {
        let timeFrame = period + 's';
        return DateTime.fromSQL(dateTo, { zone: 'UTC' }).plus({ [timeFrame]: -periodLength }).toSQL({ includeOffset: false });
    };

    getPeriodLength = (period, dateFrom, dateTo) => {
        let periodLength = 12;
        let timeFrame = period + 's';

        let i1 = DateTime.fromSQL(dateFrom, { zone: 'UTC' }),
            i2 = DateTime.fromSQL(dateTo, { zone: 'UTC' });

        let td = i2.diff(i1, timeFrame).toObject();

        periodLength = td[timeFrame];
        if (!Number.isInteger(periodLength)) {
            periodLength = Math.ceil(periodLength);
        }


        if (periodLength > 300) {
            periodLength = 300;
        }

        return periodLength;
    }

    prepareChartData = async(rawData, values = []) => {
        let barChartData = JSON.parse(JSON.stringify(barDataSample));
        let lineChartData = JSON.parse(JSON.stringify(lineDataSample));
        let lineChartValueData = JSON.parse(JSON.stringify(lineValueDataSample));
        let lChartOptions = this.state.lChartOptions;
        let lChartValueOptions = this.state.lChartValueOptions;

        if (rawData === null) {
            this.setState({barChartData, lineChartData, lineChartValueData})
        } else {
            for(var i = 0; i < rawData.length; i++){
                let label = this.prepareChartLabel(rawData[i].date_from, rawData[i].date_to);
                
                barChartData.datasets[0].data.push(rawData[i].num_of_successes);
                barChartData.datasets[1].data.push(rawData[i].num_of_failures);
                barChartData.labels.push(label);

                lineChartData.datasets[0].data.push(rawData[i].average_duration);
                lineChartData.labels.push(label);
            }

            lChartOptions.plugins["tooltip"]["callbacks"]["label"] = function(tooltipItem, data) {
                return "Average duration: " + prettyMilliseconds(tooltipItem.raw, {verbose: true});
            };

            this.setState({barChartData, lineChartData, lChartOptions});
        }

        for(i = 0; i < values.length; i++){
            let valueLabel = DateTime.fromSQL(values[i].executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ month: 'short', year: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' });

            lineChartValueData.datasets[0].data.push(Number(values[i].metric_value));
            lineChartValueData.labels.push(valueLabel);
        }

        lChartValueOptions.plugins["tooltip"]["callbacks"]["label"] = function(tooltipItem, data) {
            return tooltipItem.label + ": " + tooltipItem.raw;
        };

        this.setState({lineChartValueData, lChartValueOptions});
    };

    prepareChartLabel(dateFrom, dateTo) {
        if (this.state.period === "day") {
            return DateTime.fromSQL(dateFrom, { zone: 'UTC' }).toLocal().toLocaleString({ month: 'short', year: '2-digit', day: '2-digit' });
        } else {
            return DateTime.fromSQL(dateFrom, { zone: 'UTC' }).toLocal().toLocaleString({ month: 'short', year: '2-digit', day: '2-digit' }) + " - " + DateTime.fromSQL(dateTo, { zone: 'UTC' }).toLocal().toLocaleString({ month: 'short', year: '2-digit', day: '2-digit' });
        }
    }

    componentDidMount = async() => {
        if (this.props.item !== null) {
            let periodLength = this.getPeriodLength(
                this.state.period,
                this.state.dateFrom,
                this.state.dateTo
            );
            await this.promisedSetState({periodLength});

            await this.handleGetStats(this.props.item.uuid);
            await this.handleGetValues(this.props.item.uuid);
            await this.handleGetAggregatedStats(this.props.item.uuid);

            this.setEmptyStat();
        } else {
            this.setState({
                stats: emptyStateStats,
                barChartData: emptyStateBarChartData,
                lineChartData: emptyStateLineChartData,
                lineChartValueData: emptyStateLineChartValueData,
                aggregatedStats: emptyStateAggregatedStats,
                periodLength: 12,
            })
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    UNSAFE_componentWillReceiveProps = async(props) => {
        if (
            props.item !== null
            &&
            (
                props.dateFrom !== this.state.dateFrom
                || props.dateTo !== this.state.dateTo
                || props.item !== this.state.item
            )
        ) {
            let dF = props.dateFrom;
            if (dF) {
                let dFParts = dF.split(' ');
                if (dFParts.length === 1) {
                    dF = dF + " 00:00:00";
                }
            }

            let dT = props.dateTo;
            if (dT) {
                let dTParts = dT.split(' ');
                if (dTParts.length === 1) {
                    dT = dT + " 23:59:59";
                }
            }

            await this.promisedSetState({
                item: props.item, 
                dateFrom: dF || DateTime.now().plus({ days: -7 }).toSQLDate() + " 00:00:00",
                dateTo: dT || DateTime.now().toSQLDate() + " 23:59:59",
                barChartData: null,
                lineChartData: null,
                lineChartValueData: null,
                stats: null,
                values: null,
                aggregatedStats: null,
            });

            let periodLength = this.getPeriodLength(
                this.state.period,
                this.state.dateFrom,
                this.state.dateTo
            );
            await this.promisedSetState({periodLength});

            await this.handleGetStats(props.item.uuid);
            await this.handleGetValues(props.item.uuid);
            await this.handleGetAggregatedStats(props.item.uuid);

            this.setEmptyStat();
        }
    }

    handleGetStats = async(uuid) => {
        /*this.setState({
            stats: {
                "num_of_successes": 4,
                "num_of_failures": 2,
                "average_duration": 1501,
                "is_last_run_successful": true,
                "last_run_duration": 669,
                "last_run_executed_at": "2023-07-28 17:15:02"
            }
        });
        return;*/

        let stats = null;

        if (this.state.type === 'pipeline') {
            stats = StatService.getPipelineStats(
                uuid, 
                this.state.dateFrom, 
                this.state.dateTo
            );
        } else {
            stats = StatService.getTaskStats(
                uuid, 
                this.state.dateFrom, 
                this.state.dateTo
            );
        }

        await Promise.resolve(stats)
            .then(async(stats) => {
                if (stats.data) {
                    var items = stats.data;
                        
                    this.setState({stats: items});
                } else {
                    this.setState({stats: emptyStateStats});
                }
            })
            .catch(() => {
                this.setState({stats: emptyStateStats});
            });
    };

    handleGetValues = async(uuid) => {
        /*this.setState({
            values: [
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": 1461.454,
                    "executed_at": "2023-07-28 14:45:54"
                },
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": 1461.454,
                    "executed_at": "2023-07-28 15:15:01"
                },
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": "br",
                    "executed_at": "2023-07-28 16:15:02"
                },
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": 1461.454,
                    "executed_at": "2023-07-28 17:15:02"
                }
            ]
        });
        return;*/

        let stats = null;

        if (this.state.type === 'pipeline') {
            stats = StatService.getPipelineValues(
                uuid, 
                this.state.dateFrom, 
                this.state.dateTo
            );
        } else {
            await this.promisedSetState({values: []});
            return;
        }

        await Promise.resolve(stats)
            .then(async(stats) => {
                if (stats.data) {
                    var items = stats.data;
                        
                    this.setState({values: items});
                } else {
                    this.setState({values: emptyStateValues});
                }
            })
            .catch(() => {
                this.setState({values: emptyStateValues});
            });
    };

    handleGetAggregatedStats = async(uuid) => {
        /*await this.promisedSetState({
            aggregatedStats: [
                {
                    "date_from": "2023-07-28 00:00:00",
                    "date_to": "2023-07-28 23:59:59",
                    "num_of_successes": 4,
                    "num_of_failures": 2,
                    "average_duration": 1501
                }
            ]
        });
        this.prepareChartData(
            this.state.aggregatedStats,
            this.state.values
        );
        return;*/

        let stats = null;

        if (this.state.type === 'pipeline') {
            stats = StatService.getPipelineAggregatedStats(
                uuid, 
                this.state.dateFrom, 
                this.state.period,
                this.state.periodLength
            );
        } else {
            stats = StatService.getTaskAggregatedStats(
                uuid, 
                this.state.dateFrom,
                this.state.period,
                this.state.periodLength
            );

            /*stats = {
                "data": [
                        {
                            "date_from": "2023-11-20 00:00:00",
                            "date_to": "2023-11-20 23:59:59",
                            "num_of_successes": 6,
                            "num_of_failures": 0,
                            "average_duration": 719
                        }
                    ]
                };*/
        }

        await Promise.resolve(stats)
            .then(async(stats) => {
                if (stats.data && stats.data !== null) {
                    var items = stats.data;

                    await items.sort((a, b) => a.date_to.localeCompare(b.date_to));
                        
                    this.setState({aggregatedStats: items});
                    this.prepareChartData(items, this.state.values);
                } else {
                    this.setState({aggregatedStats: {}});
                    this.prepareChartData(null, this.state.values);
                }
            })
            .catch(() => {
                this.setState({aggregatedStats: {}});
                this.prepareChartData(null, this.state.values);
            });
    };

    onChangePeriod = async(period) => {
        let periodLength = this.getPeriodLength(
            period,
            this.state.dateFrom,
            this.state.dateTo
        );
        await this.promisedSetState({periodLength, period});

        await this.handleGetAggregatedStats(this.props.item.uuid);

        this.setEmptyStat();
    };

    setEmptyStat = () => {
        if (
            this.state.stats === null || 
            (this.state.stats.num_of_successes === 0
            && this.state.stats.num_of_failures === 0
            && this.state.stats.average_duration === 0
            && this.state.stats.last_run_duration === 0)
        ) {
            this.setState({
                barChartData: emptyStateBarChartDataZeros,
                lineChartData: emptyStateLineChartDataZeros,
                lineChartValueData: emptyStateLineChartValueDataZeros,
                aggregatedStats: emptyStateAggregatedStats,
                periodLength: 12,
            })
        }
    }

    render() {

        const {
            stats, 
            lChartOptions, 
            lChartValueOptions, 
            aggregatedStats, 
            item, 
            barChartData, 
            lineChartData, 
            lineChartValueData, 
            period, 
            values
        } = this.state;

        return (
            <>
                { stats !== null ?
                    <div className={item === null ? "emptyState": "allGood"}>
                        <Row className="pt-4 pb-5">
                            <Col sm={6} className="px-5">
                                <h3 className="mb-4">Last run</h3>
                                    {
                                        stats.num_of_successes > 0 || stats.num_of_failures > 0
                                        ? stats.is_last_run_successful === true
                                            ? <span className="statValue runSuccess">Success</span>
                                            : <span className="statValue runFailure">Failure</span>
                                        : <span className="statValue">-</span>
                                    }
                                <br/>
                                <span className="statKey">Status</span><br/><br/>

                                <span className="statValue">
                                    {stats.num_of_successes !== 0 || stats.num_of_failures !== 0 
                                        ? DateTime.fromSQL(stats.last_run_executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', millisecond: '3-digit' })
                                        : "-"
                                    }
                                </span><br/>
                                <span className="statKey">Run at</span><br/><br/>

                                <span className="statValue">
                                    {
                                        stats.num_of_successes !== 0 || stats.num_of_failures !== 0 
                                            ? stats.last_run_duration > 0
                                                ? prettyMilliseconds(stats.last_run_duration, {verbose: true})
                                                : "< 1 millisecond"
                                            : 0
                                    }
                                </span><br/>
                                <span className="statKey">Duration</span><br/><br/>
                            </Col>
                            <Col sm={6}>
                                <h3 className="mb-4">Total during the time period</h3>
                                <span className="statValue">{stats.num_of_successes}</span><br/>
                                <span className="statKey">Successful runs</span><br/><br/>

                                <span className="statValue">{stats.num_of_failures}</span><br/>
                                <span className="statKey">Failures</span><br/><br/>

                                <span className="statValue">
                                    {
                                        stats.num_of_successes !== 0 || stats.num_of_failures !== 0 
                                            ? stats.average_duration > 0
                                                ? prettyMilliseconds(stats.average_duration, {verbose: true})
                                                : "< 1 millisecond"
                                            : 0
                                    }
                                </span><br/>
                                <span className="statKey">Average duration</span><br/><br/>
                            </Col>
                        </Row>
                    </div>
                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }

                <div className={item === null ? "w-98 emptyState": "w-98"}>
                    { values !== null
                        && lineChartValueData !== null
                        ? 
                        ( values.length > 0 &&
                            <>
                                { item !== null && item.type === PIPELINE_TYPE_METRIC 
                                && <>
                                        <h3 className="mb-4">Values</h3>
                                        <div className="chartContainer mb-5">
                                            <LineChart
                                                data={lineChartValueData}
                                                options={lChartValueOptions}
                                            />
                                        </div>
                                    </>
                                }
                            </>
                        )
                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                    }
                </div>

                { item !== null &&
                <div className="text-right">
                    <Dropdown size="lg" className="mb-4 float-right">
                            <Dropdown.Toggle 
                                variant="light" 
                                id="dropdown-basic"
                                className=""
                            >
                                {period.charAt(0).toUpperCase() + period.slice(1)}
                            </Dropdown.Toggle>

                            <Dropdown.Menu className="tasks">
                                {timeFrames.map(value => (
                                    <Dropdown.Item
                                        value={value}
                                        key={value}
                                        active={value === period}
                                        onClick={(e) => this.onChangePeriod(value)}
                                    >
                                        {value.charAt(0).toUpperCase() + value.slice(1)}
                                    </Dropdown.Item>
                                ))}
                            </Dropdown.Menu>
                        </Dropdown>
                        <div className="dropdownLabel float-right">Group by: </div>
                        <div className="clearfix"></div>
                </div>
            }
                <div className={item === null ? "w-98 emptyState": "w-98"}>
                    { aggregatedStats !== null 
                        && barChartData !== null 
                        && lineChartData !== null
                        && lineChartValueData !== null
                        ? 
                        ( aggregatedStats.length > 0 &&
                            <>
                                <h3 className="mb-4">Runs</h3>
                                <div className="chartContainer">
                                    <BarChart
                                        data={barChartData}
                                    />
                                </div>
                                <h3 className="mb-4 mt-5">Average duration</h3>
                                <div className="chartContainer">
                                    <LineChart
                                        data={lineChartData}
                                        options={lChartOptions}
                                    />
                                </div>
                            </>
                        )
                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                    }
                </div>
                <div className="pb-5"></div>
            </>
        );
    }
}

export default PipelineStatisticSummary;
