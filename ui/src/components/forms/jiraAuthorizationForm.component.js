import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import { Navigate } from 'react-router-dom';

import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "../formControls/input.component";
import { required } from "../formControls/validations";

import { connect } from "react-redux";
import {
    updateJiraAuthorization
} from "../../actions/integrations";

import {
    Card,
    Col,
    Container, Dropdown,
    InputGroup,
    Row,
} from 'react-bootstrap'

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import {clearMessage} from "../../actions/message";
import log from "loglevel"

class JiraAuthorizationForm extends Component {
    constructor(props) {
        super(props);

        this.state = {
            item: this.props.item,
            name: "",
            resourceId: "",
            resourceName: "Please select",
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
        log.debug("JiraAuthorizationForm component did mount")
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            log.debug("init item", this.props.item)
            this.setState({
                item: this.props.item,
                name: this.props.item.model.name || "",
                uuid: this.props.item.model.uuid || "",
                resourceId: this.props.item.model.resource_id || "",
                justConnected: this.props.justConnected || false
            });

            var element = this.props.item.resources.find(o => o.id === this.props.item.model.resource_id);
            console.log(element)
            if (element) {
                this.setState({resourceName: element.name})
            }
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
            updateJiraAuthorization(
                this.state.item.model.uuid,
                this.state.name,
                this.state.resourceId
            )
        )
            .then(() => {
                var item = this.state.item;
                item.model.name = this.state.name;
                item.model.resource_id = this.state.resourceId
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

    onChangeJiraResource(id, name) {
        this.setState({
            resourceId: id,
            resourceName: name,
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
                                    {justConnected && !this.state.item.model.is_active && (
                                        <div className="form-group">
                                            <div className={ "alert alert-success mt-3" } role="alert">
                                                One more step. Please select your jira instance from the list below and save it to activate your jira authorization.
                                            </div>
                                        </div>
                                    )}
                                    {justConnected && this.state.item.model.is_active && (
                                        <div className="form-group">
                                            <div className={ "alert alert-success mt-3" } role="alert">
                                                Congratulations! Your jira is now connected. Feel free to rename the authorization, if you are up to.
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
                                        <InputGroup className="mb-4">
                                            <div className="float-left">
                                                <span className="inputLabel">Jiras you have access to</span><br/>
                                                <Dropdown size="lg">
                                                    <Dropdown.Toggle
                                                        variant="light"
                                                        id="dropdown-basic"
                                                    >
                                                        {this.state.resourceName}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {this.state.item.resources.map(value => (
                                                                <Dropdown.Item
                                                                    value={value.name}
                                                                    key={value.id}
                                                                    active={value.id === this.state.resourceId}
                                                                    onClick={(e) => this.onChangeJiraResource(value.id, value.name)}
                                                                >
                                                                    {value.name}
                                                                </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>
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
                                                    {message.item ? "Jira authorization updated" : message}
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

export default connect(mapStateToProps)(JiraAuthorizationForm);