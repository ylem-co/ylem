import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import { Navigate } from 'react-router-dom';

import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "../formControls/input.component";
import { required } from "../formControls/validations";

import { connect } from "react-redux";
import {
    updateSlackAuthorization
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
import IntegrationService from "../../services/integration.service";

class SlackAuthorizationForm extends Component {
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

        this.handleUpdate = this.handleUpdate.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.handleAuthorize = this.handleAuthorize.bind(this);
    };

    componentDidMount() {
        log.debug("SlackAuthorizationForm component did mount")
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            log.debug("init item", this.props.item)
            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                uuid: this.props.item.uuid || "",
                justConnected: this.props.justConnected || false
            });
        }
    }

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

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
            updateSlackAuthorization(
                this.state.item.uuid,
                this.state.name
            )
        )
            .then(() => {
                var item = this.state.item;
                item.name = this.state.name;
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

    handleAuthorize = async(e) => {
        log.debug("handle authorize")
        e.preventDefault();

        await this.promisedSetState({
            loading: true,
        });

        let link = await this.createReAuthorizationLink(this.state.uuid);
        log.debug("got an auth link: ", link)

        if (link !== null) {
            window.location.href = link;
        } else {
            log.warn("Something went off. The link was not created");

            await this.promisedSetState({
                loading: false,
            });
        }
    }

    createReAuthorizationLink = async(uuid) => {
        log.debug("handle re authorization")
        let authorizationLink = IntegrationService.reauthorizeSlackAuthorization(uuid);

        return Promise.resolve(authorizationLink)
            .then(authorizationLink => {
                if (authorizationLink.data) {
                    return authorizationLink.data.url;
                } else {
                    return null;
                }
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
                                    {justConnected && (
                                        <div className="form-group">
                                            <div className={ "alert alert-success mt-3" } role="alert">
                                                Congratulations! Your Slack has been authorized. If you want, you can rename your connection.
                                            </div>
                                        </div>
                                    )}
                                    {!this.state.item.is_active && (
                                        <div className="form-group">
                                            <div className={ "alert alert-warning mt-3" } role="alert">
                                                Your slack is either not authorized yet or needs to be re-authorized. <span onClick={this.handleAuthorize}>Please follow this link to authorize</span>.
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
                                                    {message.item ? "Slack authorization updated" : message}
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

export default connect(mapStateToProps)(SlackAuthorizationForm);
