import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Cron from "react-js-cron";
import 'react-js-cron/dist/styles.css';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import { connect } from "react-redux";
import { updatePipeline } from "../../actions/pipelines";

import Input from "../formControls/input.component";
import { isCron } from "../formControls/validations";

import { clearMessage } from "../../actions/message";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import {
    Button,
    Col,
    InputGroup,
    Row
} from 'react-bootstrap'

class ScheduleForm extends Component {
    constructor(props) {
        super(props);
        this.handleUpdate = this.handleUpdate.bind(this);
        this.onChangeSchedule = this.onChangeSchedule.bind(this);
        this.onChangeCrontab = this.onChangeCrontab.bind(this);
        this.enableScheduleEditor = this.enableScheduleEditor.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            schedule: "",
            scheduleEditorEnabled: false,
            item: this.props.item,
            loading: false,
            successful: false,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            this.setState({
                item: this.props.item,
                schedule: this.props.item.schedule || "",
                scheduleEditorEnabled: this.props.item.schedule !== null,
                parentFolder: this.props.parentFolder || null,
            });
        }
    };

    onChangeSchedule(schedule) {
        this.setState({
            schedule,
            scheduleEditorEnabled: schedule !== "",
        });
    }

    enableScheduleEditor() {
        this.setState({
            scheduleEditorEnabled: true,
        });
    }

    onChangeCrontab(e) {
        this.setState({
            schedule: e.target.value,
            scheduleEditorEnabled: e.target.value !== "",
        });
    }

    handleUpdate(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                updatePipeline(
                    this.props.item.uuid,
                    this.props.item.name, 
                    this.props.item.folder_uuid, 
                    JSON.parse(this.props.item.elements_layout), 
                    this.state.schedule
                )
            )
            .then(() => {
                this.setState({
                    loading: false,
                    successful: true,
                });

                setTimeout(() => {
                    this.props.successHandler(this.state.schedule);
                }, 1000);
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    render() {
        const { isLoggedIn, user, message, item } = this.props;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div>
                <Form
                    onSubmit={
                        item === null
                        ? this.handleCreate
                        : this.handleUpdate
                    }
                    ref={(c) => {
                        this.form = c;
                    }}
                >
                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <InputGroup className="mb-4">
                                <div className="registrationFormControl">
                                    <div className="p-3" onClick={this.enableScheduleEditor}>
                                        <Cron 
                                            value={this.state.schedule}
                                            setValue={this.onChangeSchedule}
                                            humanizeValue={false}
                                            humanizeLabels={true}
                                            leadingZero
                                            defaultPeriod={'minute'}
                                            clearButtonAction={'empty'}
                                            disabled={!this.state.scheduleEditorEnabled}
                                            allowEmpty={'always'}
                                        />
                                    </div>
                                    <span className="inputLabel">Or write it in a cronjob format</span><br/>
                                    <Input
                                        className="form-control"
                                        type="text"
                                        id="floatingSchedule"
                                        placeholder=""
                                        autoComplete="schedule"
                                        name="schedule"
                                        value={this.state.schedule}
                                        onChange={this.onChangeCrontab}
                                        validations={[isCron]}
                                    />
                                </div>
                            </InputGroup>
                        </div>
                    </InputGroup>

                    <Row>
                        <Col xs="6">
                            <Button
                                className="px-4 btn btn-primary"
                                disabled={this.state.loading}
                                type="submit"
                            >
                                {this.state.loading && (
                                    <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                )}
                                <span>Save</span>
                            </Button>
                        </Col>
                    </Row>
                    {message && (
                        <div className="form-group">
                            <div className={ this.state.successful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                {
                                    message.item 
                                        ? 
                                            (
                                                item === null
                                                ? "Folder successfully created"
                                                : "Folder successfully updated"  
                                            )
                                        : message
                                }
                            </div>
                        </div>
                    )}
                    <CheckButton
                        style={{ display: "none" }}
                            ref={(c) => {
                                this.checkBtn = c;
                            }}
                    />
                </Form>                  
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { message } = state.message;
    const { user } = state.auth;
    return {
        isLoggedIn,
        message,
        user
    };
}

export default connect(mapStateToProps)(ScheduleForm);
