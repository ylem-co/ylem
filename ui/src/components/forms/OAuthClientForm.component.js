import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import Tooltip from '@mui/material/Tooltip';
import ContentCopy from '@mui/icons-material/ContentCopy';

import FloatingLabel from "react-bootstrap/FloatingLabel";

import { connect } from "react-redux";
import { addOAuthClient } from "../../actions/OAuthClients";

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

function copyLink(copiedLink){
  navigator.clipboard.writeText(copiedLink);
}

class OAuthClientForm extends Component {
    constructor(props) {
        super(props);
        this.handleCreate = this.handleCreate.bind(this);
        this.onChangeName = this.onChangeName.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: "",
            loading: false,
            successful: false,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
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
                addOAuthClient(
                    this.state.name
                )
            )
            .then(() => {
                this.setState({
                    loading: false,
                    successful: true,
                });
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
                    onSubmit={this.handleCreate}
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

                    <Row>
                        <Col xs="6">
                            <Button
                                className="px-4 btn btn-primary"
                                disabled={this.state.loading || (message && this.state.successful)}
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
                                    message.client 
                                        ? 
                                            <div>
                                                Client successfully created<br/>
                                                Please copy the following client secret and save it securely:
                                                <div className="code withBorder">
                                                    <Row>
                                                        <Col xs={10} className="note">
                                                            {message.client.data.secret}
                                                        </Col>
                                                        <Col xs={2} className="text-right">
                                                            <Tooltip title="Click to copy to clipboard" placement="left">
                                                                <ContentCopy className="note pointer"
                                                                    onClick={() => copyLink(message.client.data.secret)}
                                                                />
                                                            </Tooltip>
                                                        </Col>
                                                    </Row>
                                                </div>
                                                You won't be able to see it again.
                                            </div>
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

export default connect(mapStateToProps)(OAuthClientForm);
