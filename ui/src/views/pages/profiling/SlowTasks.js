import React from 'react';
import { Navigate, useParams, useNavigate } from 'react-router-dom';
import {connect} from "react-redux";
import { Fade } from "react-awesome-reveal";
import { DateTime } from "luxon";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import Tooltip from '@mui/material/Tooltip';

import {SlowTasksInfo} from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import PipelineSlowTaskRuns from "../../../components/pipelines/pipelineSlowTaskRuns.component";

import {TASK_TYPE_QUERY} from "../../../services/task.service";

import {PERMISSION_LOGGED_IN, validatePermissions} from "../../../actions/pipeline";

function withParams(Component) {
  return props => <Component {...props} params={useParams()} history={useNavigate()} />;
}

class SlowTasks extends React.Component {
    constructor(props) {
        super(props);

        this.changeDatesHandler = this.changeDatesHandler.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isInfoOpen: false,
            loading: false,
            dateFrom: this.props.params.dateFrom || DateTime.now().toSQLDate() + " 00:00:00",
            dateTo: this.props.params.dateTo || DateTime.now().toSQLDate() + " 23:59:59",
            threshold: this.props.params.threshold || (localStorage.getItem('SlowTasksThreshold') || 10000),
            taskType: this.props.params.taskType || TASK_TYPE_QUERY,
        };
    }

    componentDidMount = async() => {
        document.title = 'Slow Tasks'
    }

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    changeDatesHandler = async(dateFrom, dateTo) => {
        await this.promisedSetState({
            dateFrom: dateFrom + '%2000:00:00',
            dateTo: dateTo + '%2023:59:59', 
        });
        this.props.history(
            '/slow-tasks'
            + '/' + dateFrom + '%2000:00:00' 
            + '/' + dateTo + '%2023:59:59'
            + '/' + this.state.threshold
            + '/' + this.state.taskType);
    };

    changeTaskTypeHandler = async(taskType) => {
        await this.promisedSetState({taskType});

        this.props.history(
            '/slow-tasks'
            + '/' + this.state.dateFrom
            + '/' + this.state.dateTo
            + '/' + this.state.threshold
            + '/' + taskType);
  };

  changeThresholdHandler = async(e) => {
    let threshold = e.target.value;

    if (threshold > 0) {
        await this.promisedSetState({threshold});

        await localStorage.setItem('SlowTasksThreshold', threshold);

        this.props.history(
            '/slow-tasks'
            + '/' + this.state.dateFrom 
            + '/' + this.state.dateTo
            + '/' + threshold
            + '/' + this.state.taskType);
    } else if (threshold.length === 0) {
        await this.promisedSetState({threshold});
    }
  };

    toogleInfo = async() => {
        await this.promisedSetState({
            isInfoOpen: !this.state.isInfoOpen,
        });
    }

    closeInfo = () => {
        this.setState({isInfoOpen: false});
    }

    render() {
        const { isInfoOpen, dateFrom, dateTo, threshold, taskType } = this.state;

        const { isLoggedIn, user } = this.props;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
            return <Navigate to="/login" />;
        }

        return (
            <Fade>
                <Row className="mb-3">
                    <Col sm="9">
                        <h1>Slow Tasks</h1>
                        <Tooltip title="Info" placement="right">
                            <div className="infoIcon" onClick={() => this.toogleInfo()}></div>
                        </Tooltip>
                    </Col>
                </Row>

                <PipelineSlowTaskRuns
                    dateFrom={dateFrom}
                    dateTo={dateTo}
                    threshold={threshold}
                    taskType={taskType}
                    changeTaskTypeHandler={this.changeTaskTypeHandler}
                    changeThresholdHandler={this.changeThresholdHandler}
                    changeDatesHandler={this.changeDatesHandler}
                />
            
                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={SlowTasksInfo.title}
                    content={SlowTasksInfo.content}
                />
            </Fade>
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

export default connect(mapStateToProps)(withParams(SlowTasks));
