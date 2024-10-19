import React, {Component} from "react";
import {Link, Navigate, useNavigate} from 'react-router-dom';
import { Fade } from "react-awesome-reveal";

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "./formControls/input.component";
import {required} from "./formControls/validations";

import {connect} from "react-redux";
import {login} from "../actions/auth";

import {PERMISSION_LOGGED_OUT, validatePermissions} from "../actions/pipeline";

import {Button, Card, CardGroup, Col, Container, InputGroup, Row} from 'react-bootstrap'

import AuthService from "../services/auth.service";

function withParams(Component) {
  return props => <Component {...props} history={useNavigate()} />;
}

class Login extends Component {
    constructor(props) {
        super(props);
        this.handleLogin = this.handleLogin.bind(this);
        this.onChangeEmail = this.onChangeEmail.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);
        this.handleSignInWithGoogle = this.handleSignInWithGoogle.bind(this);
        this.handleGetSignInWithGoogleAvailability = this.handleGetSignInWithGoogleAvailability.bind(this);

        this.state = {
            email: "",
            password: "",
            loading: false,
            isGoogleAuthAvailable: null,
            isDarkThemeEnabled: localStorage.getItem('darkTheme') !== "false",
        };
    }

    componentDidMount() {
        document.title = 'Login';

        this.handleGetSignInWithGoogleAvailability();
    };

    onChangeEmail(e) {
        this.setState({
            email: e.target.value,
        });
    }

    onChangePassword(e) {
        this.setState({
            password: e.target.value,
        });
    }

    handleLogin = async(e) => {
        e.preventDefault();

        this.setState({
            loading: true,
        });

        this.form.validateAll();

        const { dispatch, history } = this.props;

        if (this.checkBtn.context._errors.length === 0) {

            dispatch(login(this.state.email, this.state.password))
                .then(() => {
                    history("/dashboard");
                })
                .catch(() => {
                    this.setState({
                        loading: false
                    });
                });
        } else {
            this.setState({
                loading: false,
            });
        }
    }

    handleSignInWithGoogle(e) {
        e.preventDefault()
        let link = AuthService.getSignInWithGoogleRedirectUrl()

        Promise.resolve(link)
            .then(link => {
                window.location = link.data.url
            })
            .catch(() => {
            });
    }

    handleGetSignInWithGoogleAvailability() {
        let isAvailable = AuthService.getSignInWithGoogleAvailability()

        Promise.resolve(isAvailable)
            .then(isAvailable => {
                this.setState({
                    isGoogleAuthAvailable: true,
                });
            })
            .catch(() => {
                this.setState({
                    isGoogleAuthAvailable: false,
                });
            });
    }

    render() {
        const { isLoggedIn, user, message } = this.props;

        const { isDarkThemeEnabled, isGoogleAuthAvailable } = this.state;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/dashboard" />
        }

        return (
        <div className="align-items-center visualBg">
            <Fade>
                <Container>
                    <Row className="justify-content-center">
                        <Col md="8" sm="12" className="text-center pt-5">
                        {isDarkThemeEnabled
                            ? <img src="/images/logo2-dark.png" width="158px" alt="Ylem"/>
                            : <img src="/images/logo2.png" width="158px" alt="Ylem"/>
                        }
                        </Col>
                    </Row>
                    <Row className="justify-content-center pt-5">
                        <Col lg="10" md="12" xl="8">
                            <CardGroup>
                                <Card className="p-4 onboardingCard">
                                    <Card.Body>
                                        <Form
                                            onSubmit={this.handleLogin}
                                            ref={(c) => {
                                                this.form = c;
                                            }}
                                        >
                                            <h2 className="mb-4">Sign in</h2>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingEmail" label="Email">
                                                <Input
                                                    className="form-control form-control-lg"
                                                    type="text"
                                                    placeholder="Email"
                                                    id="floatingEmail"
                                                    autoComplete="email"
                                                    name="email"
                                                    value={this.state.email}
                                                    onChange={this.onChangeEmail}
                                                    autoFocus
                                                    validations={[required]}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingPassword" label="Password">
                                                <Input
                                                    className="form-control form-control-lg"
                                                    type="password"
                                                    id="floatingPassword"
                                                    placeholder="Password"
                                                    autoComplete="current-password"
                                                    value={this.state.password}
                                                    onChange={this.onChangePassword}
                                                    validations={[required]}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                            <Row>
                                                <Col xs="5">
                                                    <Button
                                                        className="px-4 btn btn-primary"
                                                        disabled={this.state.loading}
                                                        type="submit"
                                                    >
                                                        {this.state.loading && (
                                                            <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                        )}
                                                        <span>Login</span>
                                                    </Button>
                                                </Col>
                                                <Col xs="7">
                                                    { isGoogleAuthAvailable === true
                                                        &&
                                                        <button className="google-sign-in-button btn-secondary btn"
                                                                onClick={this.handleSignInWithGoogle}>
                                                            Continue with
                                                            <span className="google-sign-in-icon"></span>
                                                        </button>
                                                    }
                                                </Col>
                                            </Row>
                                            {message && (
                                                <div className="form-group">
                                                    <div className="alert alert-danger mt-3" role="alert">
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
                                <Card className="onboardingCard py-5" style={{ width: '44%' }}>
                                    <Card.Body className="text-center">
                                        <div>
                                            <h2>Sign up</h2>
                                            <h4>Do not have an account yet?</h4>
                                            <Link to="/register">
                                                <Button color="secondary" variant="secondary" className="mt-3" active tabIndex={-1}>Create</Button>
                                            </Link>
                                        </div>
                                    </Card.Body>
                                </Card>
                            </CardGroup>
                        </Col>
                    </Row>
                </Container>
            </Fade>
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

export default connect(mapStateToProps)(withParams(Login));
