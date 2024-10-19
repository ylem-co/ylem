import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import { Navigate, Link, useParams, useNavigate} from 'react-router-dom';

import Spinner from "react-bootstrap/Spinner";
import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "./formControls/input.component";
import { required, isEqual } from "./formControls/validations";

import VisibilityOutlined from '@mui/icons-material/VisibilityOutlined';
import VisibilityOffOutlined from '@mui/icons-material/VisibilityOffOutlined';

import Tooltip from '@mui/material/Tooltip';

import { connect } from "react-redux";
import { register } from "../actions/auth";

import {
    Card,
    Col,
    Container,
    InputGroup,
    Row
} from 'react-bootstrap'

import InvitationService from "../services/invitation.service";

import { validatePermissions, PERMISSION_LOGGED_OUT }   from "../actions/pipeline";

function withParams(Component) {
  return props => <Component {...props} params={useParams()} history={useNavigate()} />;
}

class Register extends Component {
    constructor(props) {
        super(props);
        this.handleRegister = this.handleRegister.bind(this);
        this.onChangeFirstName = this.onChangeFirstName.bind(this);
        this.onChangeLastName = this.onChangeLastName.bind(this);
        this.onChangeEmail = this.onChangeEmail.bind(this);
        this.onChangePhone = this.onChangePhone.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);
        this.onChangeConfirmPassword = this.onChangeConfirmPassword.bind(this);
        this.onChangeTermsAndConditions = this.onChangeTermsAndConditions.bind(this);
        this.onChangeNewsAndUpdates = this.onChangeNewsAndUpdates.bind(this);

        this.state = {
            firstName: "",
            lastName: "",
            phone: "",
            email: "",
            password: "",
            passwordType: "password",
            confirmPassword: "",
            confirmPasswordType: "password",
            termsAndConditions: "",
            newsAndUpdates: false,
            successful: false,
            loading: false,
            invitationKey: this.props.params.key || null,
            isKeyValid: null,
            isDarkThemeEnabled: localStorage.getItem('darkTheme') !== "false",
        };
    }

    componentDidMount() {
        document.title = 'Registration by invitation'
        var response = InvitationService.validateInvitationKey(this.props.params.key);

        Promise.resolve(response)
            .then(response => {
                this.setState({
                    isKeyValid: true,
                });
            })
            .catch(() => {
                this.props.history('/404');
            });
    };

    onChangeFirstName(e) {
        this.setState({
            firstName: e.target.value,
        });
    }

    onChangeLastName(e) {
        this.setState({
            lastName: e.target.value,
        });
    }

    onChangeEmail(e) {
        this.setState({
            email: e.target.value,
        });
    }

    onChangePhone(e) {
        this.setState({
            phone: e.target.value,
        });
    }

    onChangePassword(e) {
        this.setState({
            password: e.target.value,
        });
    }

    onChangeConfirmPassword(e) {
        this.setState({
            confirmPassword: e.target.value,
        });
    }

    onChangeTermsAndConditions(e) {
        this.setState({
            termsAndConditions: e.target.value,
        });
    }

    onChangeNewsAndUpdates(e) {
        this.setState({
            newsAndUpdates: e.target.checked,
        });
    }

    handleEyeClick = () => this.setState(({passwordType}) => ({
        passwordType: passwordType === 'text' ? 'password' : 'text'
    }));

    handleConfirmEyeClick = () => this.setState(({confirmPasswordType}) => ({
        confirmPasswordType: confirmPasswordType === 'text' ? 'password' : 'text'
    }));

    handleRegister = async(e) => {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: true,
        });

        this.form.validateAll();

        if (this.checkBtn.context._errors.length === 0) {
            this.props
                .dispatch(
                    register(
                        this.state.firstName,
                        this.state.lastName,
                        this.state.email, 
                        this.state.password, 
                        this.state.confirmPassword,
                        this.state.phone, 
                        null,
                        this.state.invitationKey
                    )
                )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: false,
                    });
                    setTimeout(() => {
                        this.props.history('/login')
                    },2000);
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: false,
                    });
                });
        } else {
            this.setState({
                loading: false,
            });
        }
    }

    render() {
        const { isLoggedIn, message, user } = this.props;

        const { isDarkThemeEnabled } = this.state;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/dashboard" />;
        }

        return (
            <div className="align-items-center visualBg">
                <Link to="/" className="floatingLogo">
                    {isDarkThemeEnabled
                        ? <img src="/images/logo2-dark.png" width="150px" alt="Ylem"/>
                        : <img src="/images/logo2.png" width="150px" alt="Ylem"/>
                    }
                </Link>
                <Container>
                    <Row className="justify-content-center pt-5">
                        { this.state.isKeyValid === true ?
                        <Col md="9" lg="7" xl="6">
                            <h2 className="alternative text-center mb-3">Create your Ylem account</h2>
                            <Card className="onboardingCard mb-5">
                                <Card.Body className="p-4">
                                    <Form
                                        onSubmit={this.handleRegister}
                                        ref={(c) => {
                                            this.form = c;
                                        }}
                                    >
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                            <FloatingLabel controlId="floatingFirstName" label="First name">
                                            <Input
                                                className="form-control form-control-lg"
                                                id="floatingFirstName"
                                                type="text"
                                                placeholder="First name"
                                                autoComplete="firstName"
                                                name="firstName"
                                                value={this.state.firstNamev}
                                                onChange={this.onChangeFirstName}
                                                autoFocus
                                                validations={[required]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                            <FloatingLabel controlId="floatingLastName" label="Last name">
                                            <Input
                                                className="form-control form-control-lg"
                                                id="floatingLastName"
                                                type="text"
                                                placeholder="Last name"
                                                autoComplete="lastName"
                                                name="lastName"
                                                value={this.state.lastName}
                                                onChange={this.onChangeLastName}
                                                validations={[required]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                            <FloatingLabel controlId="floatingEmail" label="Email">
                                            <Input
                                                type="text"
                                                id="floatingEmail"
                                                placeholder="Email"
                                                autoComplete="email"
                                                className="form-control form-control-lg"
                                                name="email"
                                                value={this.state.email}
                                                onChange={this.onChangeEmail}
                                                validations={[required]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                            <FloatingLabel controlId="floatingPhone" label="Phone (optional)">
                                            <Input
                                                type="text"
                                                id="floatingPhone"
                                                placeholder="Phone (optional)"
                                                autoComplete="phone"
                                                className="form-control form-control-lg"
                                                name="phone"
                                                value={this.state.phone}
                                                onChange={this.onChangePhone}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                            <FloatingLabel controlId="floatingPassword" label="Password">
                                            <Input
                                                type={this.state.passwordType}
                                                id="floatingPassword"
                                                placeholder="Password"
                                                autoComplete="new-password"
                                                className="form-control form-control-lg"
                                                name="password"
                                                value={this.state.password}
                                                onChange={this.onChangePassword}
                                                validations={[required, isEqual]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                            <span
                                                onClick={this.handleEyeClick}
                                                className="eye"
                                            >
                                                {
                                                    this.state.passwordType === 'text' 
                                                    ? <Tooltip title="Hide" placement="right"><VisibilityOffOutlined/></Tooltip>
                                                     : <Tooltip title="Show" placement="right"><VisibilityOutlined/></Tooltip>
                                                 }
                                            </span>
                                        </InputGroup>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                            <FloatingLabel controlId="floatingConfirmPassword" label="Confirm password">
                                            <Input
                                                type={this.state.confirmPasswordType}
                                                id="floatingConfirmPassword"
                                                placeholder="Confirm password"
                                                autoComplete="confirm-password"
                                                className="form-control form-control-lg"
                                                name="confirmPassword"
                                                value={this.state.confirmPassword}
                                                onChange={this.onChangeConfirmPassword}
                                                validations={[required, isEqual]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                            <span
                                                onClick={this.handleConfirmEyeClick}
                                                className="eye"
                                            >
                                                {
                                                    this.state.confirmPasswordType === 'text' 
                                                    ? <Tooltip title="Hide" placement="right"><VisibilityOffOutlined/></Tooltip>
                                                     : <Tooltip title="Show" placement="right"><VisibilityOutlined/></Tooltip>
                                                }
                                            </span>
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
                                                    <span>Create</span>
                                                </button>
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
                                </Card.Body>
                            </Card>
                        </Col>
                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                    }
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

export default connect(mapStateToProps)(withParams(Register));
