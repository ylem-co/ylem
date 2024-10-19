import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";
import Spinner from "react-bootstrap/Spinner";

import { connect } from "react-redux";
import { updateTask } from "../../../actions/tasks";

import Input from "../../formControls/input.component";
import InputChips from "../../formControls/inputChips.component";

import { clearMessage } from "../../../actions/message";

import { TASK_SEVERITY_MEDIUM } from "../../../services/task.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    InputGroup,
    Row
} from 'react-bootstrap'

class MergeForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeFieldNames = this.onChangeFieldNames.bind(this);

        this.state = {
            name: "",
            severity: TASK_SEVERITY_MEDIUM,
            fieldNames: null,
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
        	this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
                fieldNames: 
                	this.props.item.implementation.field_names 
                		? this.props.item.implementation.field_names.split(",")
                		: [],
            });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeFieldNames = (fieldNames) => {
        this.setState({fieldNames});
    };

    handleSubmit(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            let name = this.state.name;
            if (name === "") {
                name = "Merge";
            }
            dispatch(
                updateTask(
                    this.state.item.uuid, 
                    this.state.item.pipeline_uuid,
                    name,
                    this.state.severity,
                    this.state.item.type,
                    {
                        "field_names": this.state.fieldNames.join(","),
                    }
                )
            )
            .then(() => {
                var item = this.state.item;
                item.name = name;
                this.setState({
                    loading: false,
                    successful: true,
                    item,
                });

                this.props.successHandler(item);
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
        const { isLoggedIn, user, message } = this.props;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div>
                <Form
                    onSubmit={this.handleSubmit}
                    ref={(c) => {
                        this.form = c;
                    }}
                >
                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <FloatingLabel controlId="floatingName" label="Name (optional)">
                                <Input
                                    className="form-control form-control-lg"
                                    type="text"
                                    id="floatingName"
                                    placeholder="Name (optional)"
                                    autoComplete="name"
                                    name="name"
                                    value={this.state.name}
                                    onChange={this.onChangeName}
                                />
                            </FloatingLabel>
                        </div>
                    </InputGroup>

                    { this.state.fieldNames !== null ?
	                    <InputGroup className="mb-4">
	                        <div className="registrationFormControl">
	                            <span className="inputLabel">Key field names (optional)</span><br/><br/>
	                            <InputChips
		                        	changesHandler={this.onChangeFieldNames}
		                        	items={this.state.fieldNames}
		                        />
	                            <div className="inputTip">
                                    If two rows have the same values in all fields enumerated in this parameter, they will be merged.<br/><br/>
                                    E.g. <br/>
                                    If you add a field "id" — the block will merge rows with same IDs, <br/>
                                    If you add fields "id" and "name" — rows with the same ID and name will be merged.<br/><br/> 
                                    If you leave key field names empty, all rows from all inputs will be merged.
                                </div>
	                        </div>
	                    </InputGroup>
	                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
	                }

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
                                {message}
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

export default connect(mapStateToProps)(MergeForm);

