import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import { Navigate } from 'react-router-dom';

import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "../formControls/input.component";
import { required } from "../formControls/validations";

import { connect } from "react-redux";
import {
    updateSalesforceAuthorization,
} from "../../actions/integrations";

import {
    Card,
    Col,
    Container,
    InputGroup,
    Row,
} from 'react-bootstrap'

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import {clearMessage} from "../../actions/message";
import log from "loglevel"

class SalesforceAuthorizationForm extends Component {
    constructor(props) {
        super(props);

        this.state = {
            item: this.props.item,
            name: "",
            uuid: "",
            isInProgress: false,
            successful: null,
            justConnected: false,
        };

        console.log(this.props.item)

        this.handleUpdate = this.handleUpdate.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
    };

    componentDidMount() {
        log.debug("SalesforceAuthorizationForm component did mount")
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            log.debug("init item", this.props.item)
            this.setState({
                item: this.props.item,
                name: this.props.item.model.name || "",
                uuid: this.props.item.model.uuid || "",
                justConnected: this.props.justConnected || false
            });
        }
    }

    handleUpdate(e) {
        log.debug("handle update", this.state)

        e.preventDefault();

        this.form.validateAll();

        if (this.checkBtn.context._errors.length !== 0) {
            log.debug("form is not valid")

            return
        }

        this.setState({
            loading: true,
            successful: false,
        });

        const { dispatch } = this.props;

        dispatch(
            updateSalesforceAuthorization(
                this.state.item.model.uuid,
                this.state.name,
            )
        )
            .then(() => {
                var item = this.state.item;
                item.model.name = this.state.name;
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
    }

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    render() {
        const { isLoggedIn, message, user } = this.props;

        const { justConnected } = this.state;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div className="align-items-center">
                <Container>
                    <Row className="justify-content-center">
                        <Col md="10" lg="8" xl="7">
                        {
                            <Card className="onboardingCard noBorder mb-5">
                                <Card.Body className="p-4">
                                    {justConnected && this.state.item.model.is_active && (
                                        <div className="form-group">
                                            <div className={ "alert alert-success mt-3" } role="alert">
                                                Congratulations! Your Salesforce is now connected. Feel free to rename the authorization, if you are up to.
                                            </div>
                                        </div>
                                    )}
                                    <Form
                                        onSubmit={this.handleUpdate}
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
                                                autoFocus
                                                validations={[required]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>
                                        <Row className="pt-3">
                                            <Col xs="12">
                                                <button
                                                    className="px-4 btn btn-primary"
                                                    disabled={this.state.loading}
                                                    type="submit"
                                                >
                                                    {this.state.loading && (
                                                        <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                    )}
                                                    <span>Save</span>
                                                </button>
                                            </Col>
                                        </Row>
                                        {message && (
                                            <div className="form-group">
                                                <div className={ this.state.successful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                                    {message.item ? "Salesforce authorization updated" : message}
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
                                </Card.Body>
                            </Card>
                        }
                        </Col>
                    </Row>
                </Container>
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    const { message } = state.message;
    return {
        isLoggedIn,
        message,
        user,
    };
}

export default connect(mapStateToProps)(SalesforceAuthorizationForm);
