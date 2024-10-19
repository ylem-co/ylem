import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";

import { connect } from "react-redux";
import { updateEnvVariable, addEnvVariable } from "../../actions/envVariables";

import Input from "../formControls/input.component";
import { required } from "../formControls/validations";

import { clearMessage } from "../../actions/message";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import {
    Button,
    Col,
    InputGroup,
    Row
} from 'react-bootstrap'

class EnvVariableForm extends Component {
    constructor(props) {
        super(props);
        this.handleCreate = this.handleCreate.bind(this);
        this.handleUpdate = this.handleUpdate.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeValue = this.onChangeValue.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: "",
            value: "",
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
                name: this.props.item.name || "",
                value: this.props.item.value || null,
            });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeValue(e) {
        this.setState({
            value: e.target.value,
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
                updateEnvVariable(
                    this.state.item.uuid,
                    this.state.name,
                    this.state.value
                )
            )
            .then(() => {
                var item = this.state.item;
                item.name = this.state.name;
                item.value = this.state.value;
                this.setState({
                    loading: false,
                    successful: true,
                    item,
                });

                setTimeout(() => {
                    this.props.successHandler();
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

    handleCreate(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                addEnvVariable(
                    this.state.name,
                    this.state.value,
                    this.state.organization.uuid
                )
            )
            .then(() => {
                this.setState({
                    loading: false,
                    successful: true,
                });

                setTimeout(() => {
                    this.props.successHandler();
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
                            <FloatingLabel controlId="floatingName" label="Name">
                                <Input
                                    className="form-control form-control-lg"
                                    type="text"
                                    id="floatingName"
                                    placeholder="Name"
                                    autoComplete="name"
                                    name="name"
                                    value={this.state.name}
                                    onChange={this.onChangeName}
                                    validations={[required]}
                                />
                            </FloatingLabel>
                        </div>
                    </InputGroup>

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <FloatingLabel controlId="floatingValue" label="Value">
                                <Input
                                    className="form-control form-control-lg"
                                    type="text"
                                    id="floatingValue"
                                    placeholder="Value"
                                    autoComplete="value"
                                    name="value"
                                    value={this.state.value}
                                    onChange={this.onChangeValue}
                                    validations={[required]}
                                />
                            </FloatingLabel>
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
                                                ? "Environment variable successfully created"
                                                : "Environment variable successfully updated"  
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

export default connect(mapStateToProps)(EnvVariableForm);