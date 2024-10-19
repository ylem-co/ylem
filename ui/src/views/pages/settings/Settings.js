import React from 'react';
import { Navigate } from 'react-router-dom';
import {connect} from "react-redux";
import {logout} from "../../../actions/auth";
import { Fade } from "react-awesome-reveal";

import VisibilityOutlined from '@mui/icons-material/VisibilityOutlined';
import VisibilityOffOutlined from '@mui/icons-material/VisibilityOffOutlined';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import InputGroup from "react-bootstrap/InputGroup";
import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "../../../components/formControls/input.component";
import { required, isEqual } from "../../../components/formControls/validations";

import Tooltip from '@mui/material/Tooltip';

import {SettingsInfo} from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";

import {PERMISSION_LOGGED_IN, validatePermissions} from "../../../actions/pipeline";

import { ROLE_ORGANIZATION_ADMIN } from "../../../actions/roles";

import { updateOrganization, updatePassword, updateMe } from "../../../actions/settings";

class Settings extends React.Component {

    constructor(props) {
        super(props);

        this.handleUpdateUser = this.handleUpdateUser.bind(this);
        this.handleUpdatePassword = this.handleUpdatePassword.bind(this);
        this.handleUpdateOrganization = this.handleUpdateOrganization.bind(this);
        this.onChangeFirstName = this.onChangeFirstName.bind(this);
        this.onChangeLastName = this.onChangeLastName.bind(this);
        this.onChangePhone = this.onChangePhone.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);
        this.onChangeConfirmPassword = this.onChangeConfirmPassword.bind(this);
        this.onChangeOrganizationName = this.onChangeOrganizationName.bind(this);

        var user = JSON.parse(localStorage.getItem('user'));
        var organization = JSON.parse(localStorage.getItem('organization'));

        this.state = {
            organization: organization,
            email: user.email,
            firstName: user.first_name,
            lastName: user.last_name,
            phone: user.phone,
            organizationName: organization.name,
            user: user,
            password: "",
            passwordType: "password",
            confirmPassword: "",
            confirmPasswordType: "password",
            loading: null,
            successful: false,
            activeForm: null,
            isInfoOpen: false,
        };
    }

    componentDidMount() {
        document.title = 'Settings'
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    logOut() {
        this.props.dispatch(logout());
    }

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

    onChangeOrganizationName(e) {
        this.setState({
            organizationName: e.target.value,
        });
    }

    handleEyeClick = () => this.setState(({passwordType}) => ({
        passwordType: passwordType === 'text' ? 'password' : 'text'
    }));

    handleConfirmEyeClick = () => this.setState(({confirmPasswordType}) => ({
        confirmPasswordType: confirmPasswordType === 'text' ? 'password' : 'text'
    }));

    handleUpdateOrganization(e) {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: 'button-org',
            activeForm: null,
        });

        this.form.validateAll();

        if (this.checkBtnOrg.context._errors.length === 0) {
            this.setState({
                activeForm: 'org',
            });
            this.props
                .dispatch(
                    updateOrganization(this.state.organization.uuid, this.state.organizationName)
                )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: null,
                    });
                    var org = this.state.organization;
                    org.name = this.state.organizationName;
                    this.setState({
                        successful: true,
                        loading: null,
                    });

                    localStorage.setItem('organization', JSON.stringify(org));
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: null,
                    });
                });
        } else {
            this.setState({
                loading: null,
            });
        }
    }

    handleUpdatePassword(e) {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: 'button-pwd',
            activeForm: null,
        });

        this.form.validateAll();

        if (this.checkBtnPwd.context._errors.length === 0) {
            this.setState({
                activeForm: 'pwd',
            });
            this.props
                .dispatch(
                    updatePassword(this.state.user.uuid, this.state.password, this.state.confirmPassword)
                )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: null,
                    });
                    setTimeout(() => {
                        this.logOut()
                    },3000);
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: null,
                    });
                });
        } else {
            this.setState({
                loading: null,
            });
        }
    }

    handleUpdateUser(e) {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: 'button-upd-user',
            activeForm: null,
        });

        this.form.validateAll();

        if (this.checkBtnUsr.context._errors.length === 0) {
            this.setState({
                activeForm: 'upd-user',
            });
            this.props
                .dispatch(
                    updateMe(this.state.user.uuid, this.state.firstName, this.state.lastName, this.state.email, this.state.phone)
                )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: null,
                    });

                    var user = this.state.user;
                    user.email = this.state.email;
                    user.first_name = this.state.firstName;
                    user.last_name = this.state.lastName;
                    user.phone = this.state.phone;

                    this.setState({user})
                    localStorage.setItem('user', JSON.stringify(user));
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: null,
                    });
                });
        } else {
            this.setState({
                loading: null,
            });
        }
    }

    toogleInfo = async() => {
        await this.promisedSetState({
            isInfoOpen: !this.state.isInfoOpen,
        });
    };

    closeInfo = () => {
        this.setState({isInfoOpen: false});
    };

    render() {
        const { isInfoOpen } = this.state;

        const { isLoggedIn, message, user } = this.props;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
            return <Navigate to="/login" />;
        }

        return (
            <Fade>
                <div>
                    <Row className="mb-3">
                        <Col sm="12">
                            <h1>Settings</h1>
                            <Tooltip title="Info" placement="right">
                                <div className="infoIcon" onClick={() => this.toogleInfo()}></div>
                            </Tooltip>
                        </Col>
                    </Row>
                </div>     
                        <Row className="clearfix">
                            <Col sm="6">
                                <Card className="mr-4 mt-3 withHeader">
                                    <Card.Header>
                                        Personal settings
                                    </Card.Header>
                                    <Card.Body>
                                        <Form
                                            onSubmit={this.handleUpdateUser}
                                            ref={(c) => {
                                                this.form = c;
                                            }}
                                        >
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingFirstName" label="First name">
                                                <Input
                                                    className="form-control"
                                                    id="floatingFirstName"
                                                    type="text"
                                                    placeholder="First name"
                                                    autoComplete="firstName"
                                                    name="firstName"
                                                    value={this.state.firstName}
                                                    onChange={(e) => this.onChangeFirstName(e)}
                                                    validations={[required]}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingLastName" label="Last name">
                                                <Input
                                                    className="form-control"
                                                    id="floatingLastName"
                                                    type="text"
                                                    placeholder="Last name"
                                                    autoComplete="lastName"
                                                    name="lastName"
                                                    value={this.state.lastName}
                                                    onChange={(e) => this.onChangeLastName(e)}
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
                                                    className="form-control"
                                                    name="email"
                                                    value={this.state.email}
                                                    readOnly={"readOnly"}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingPhone" label="Phone (optional)">
                                                <Input
                                                    className="form-control"
                                                    id="floatingPhone"
                                                    type="text"
                                                    placeholder="Phone (optional)"
                                                    autoComplete="phone"
                                                    name="phone"
                                                    value={this.state.phone}
                                                    onChange={(e) => this.onChangePhone(e)}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                            <Row>
                                                <Col xs="9">
                                                    <button
                                                        id="button-upd-user"
                                                        className="px-4 btn btn-primary"
                                                        disabled={this.state.loading !== null}
                                                    >
                                                        {this.state.loading === 'button-upd-user' && (
                                                            <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                        )}
                                                        <span>Save</span>
                                                    </button>
                                                </Col>
                                            </Row>
                                            {(message && this.state.loading === null && this.state.activeForm === 'upd-user') && (
                                                <div className="form-group mt-4">
                                                    <div className={ this.state.successful ? "alert alert-success" : "alert alert-danger" } role="alert">
                                                        {message}
                                                    </div>
                                                </div>
                                            )}
                                            <CheckButton
                                                style={{ display: "none" }}
                                                ref={(c) => {
                                                    this.checkBtnUsr = c;
                                                }}
                                            />
                                        </Form>
                                    </Card.Body>
                                </Card>
                            </Col>
                            <Col sm="6">
                                <Card className="mr-4 mt-3 withHeader">
                                    <Card.Header>
                                        Change password
                                    </Card.Header>
                                    <Card.Body>
                                        <Form
                                            onSubmit={this.handleUpdatePassword}
                                            ref={(c) => {
                                                this.form = c;
                                            }}
                                        >
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingNewPassword" label="New password">
                                                <Input
                                                    type={this.state.passwordType}
                                                    id="floatingNewPassword"
                                                    placeholder="New password"
                                                    autoComplete="new-password"
                                                    className="form-control"
                                                    name="password"
                                                    value={this.state.password}
                                                    onChange={(e) => this.onChangePassword(e)}
                                                    validations={[required, isEqual]}
                                                />
                                                </FloatingLabel>
                                                </div>
                                                <span
                                                    onClick={this.handleEyeClick}
                                                    className="eye smallEye"
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
                                                <FloatingLabel controlId="floatingConfirmNewPassword" label="Confirm new password">
                                                <Input
                                                    type={this.state.confirmPasswordType}
                                                    id="floatingConfirmNewPassword"
                                                    placeholder="Confirm new password"
                                                    autoComplete="confirm-password"
                                                    className="form-control"
                                                    name="confirmPassword"
                                                    value={this.state.confirmPassword}
                                                    onChange={(e) => this.onChangeConfirmPassword(e)}
                                                    validations={[required, isEqual]}
                                                />
                                                </FloatingLabel>
                                                </div>
                                                <span
                                                    onClick={this.handleConfirmEyeClick}
                                                    className="eye smallEye"
                                                >
                                                {
                                                    this.state.confirmPasswordType === 'text' 
                                                    ? <Tooltip title="Hide" placement="right"><VisibilityOffOutlined/></Tooltip>
                                                     : <Tooltip title="Show" placement="right"><VisibilityOutlined/></Tooltip>
                                                 }
                                            </span>
                                            </InputGroup>
                                            <Row>
                                                <Col xs="9">
                                                    <button
                                                        id="button-pwd"
                                                        className="px-4 btn btn-primary"
                                                        disabled={this.state.loading !== null}
                                                    >
                                                        {this.state.loading === 'button-pwd' && (
                                                            <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                        )}
                                                        <span>Save</span>
                                                    </button>
                                                </Col>
                                            </Row>
                                            {(message && this.state.loading === null && this.state.activeForm === 'pwd') && (
                                                <div className="form-group mt-4">
                                                    <div className={ this.state.successful ? "alert alert-success" : "alert alert-danger" } role="alert">
                                                        {message}
                                                    </div>
                                                </div>
                                            )}
                                            <CheckButton
                                                style={{ display: "none" }}
                                                ref={(c) => {
                                                    this.checkBtnPwd = c;
                                                }}
                                            />
                                        </Form>
                                    </Card.Body>
                                </Card>
                                {
                                    this.state.organization !== null 
                                    && user !== null 
                                    && user.roles 
                                    && user.roles.includes(ROLE_ORGANIZATION_ADMIN) 
                                    &&
                                            <Card className="mt-4 withHeader">
                                                <Card.Header>
                                                    Organization settings
                                                </Card.Header>
                                                <Card.Body>
                                                    <Form
                                                        onSubmit={this.handleUpdateOrganization}
                                                        ref={(c) => {
                                                            this.form = c;
                                                        }}
                                                    >
                                                        <InputGroup className="mb-4">
                                                            <div className="registrationFormControl">
                                                            <FloatingLabel controlId="floatingOrganizationName" label="Organization name">
                                                            <Input
                                                                type="text"
                                                                id="floatingOrganizationName"
                                                                placeholder="Organization name"
                                                                autoComplete="organizationName"
                                                                className="form-control"
                                                                name="organizationName"
                                                                value={this.state.organizationName}
                                                                onChange={(e) => this.onChangeOrganizationName(e)}
                                                                validations={[required]}
                                                            />
                                                            </FloatingLabel>
                                                            </div>
                                                        </InputGroup>
                                                        <Row>
                                                            <Col xs="9">
                                                                <button
                                                                    id="button-org"
                                                                    className="px-4 btn btn-primary"
                                                                    disabled={this.state.loading !== null}
                                                                >
                                                                    {this.state.loading === 'button-org' && (
                                                                        <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                                    )}
                                                                    <span>Save</span>
                                                                </button>
                                                            </Col>
                                                        </Row>
                                                        {(message && this.state.loading === null && this.state.activeForm === 'org') && (
                                                            <div className="form-group mt-4">
                                                                <div className={ this.state.successful ? "alert alert-success" : "alert alert-danger" } role="alert">
                                                                    {message}
                                                                </div>
                                                            </div>
                                                        )}
                                                        <CheckButton
                                                            style={{ display: "none" }}
                                                            ref={(c) => {
                                                                this.checkBtnOrg = c;
                                                            }}
                                                        />
                                                    </Form>
                                                </Card.Body>
                                            </Card>
                                }
                            </Col>
                        </Row>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={SettingsInfo.title}
                    content={SettingsInfo.content}
                />
            </Fade>
        )
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { message } = state.message;
    const { user } = state.auth;
    return {
        isLoggedIn,
        message,
        user,
    };
}

export default connect(mapStateToProps)(Settings);
