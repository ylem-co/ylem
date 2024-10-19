import React, { Component } from "react";
import { Link, useParams, useNavigate } from 'react-router-dom';

import { connect } from "react-redux";

import {
    Card,
    Col,
    Container,
    Row,
    Spinner
} from 'react-bootstrap';

import UserService from "../services/user.service";

function withParams(Component) {
  return props => <Component {...props} params={useParams()} history={useNavigate()} />;
}

class ConfirmEmail extends Component {
    constructor(props) {
        super(props);
        this.handleConfirmEmail = this.handleConfirmEmail.bind(this);

        this.state = {
            loading: false,
            successful: null,
            key: this.props.params.key || null,
            isDarkThemeEnabled: localStorage.getItem('darkTheme') !== "false",
        };
    }

    componentDidMount() {
        document.title = 'Confirm Email';
        if (this.state.key !== null) {
            this.handleConfirmEmail(this.state.key);
        } else {
            this.setState({
                successful: false,
            });
        }
    };

    handleConfirmEmail(e) {
        this.setState({
            loading: true,
        });

        let confirmed = UserService.confirmEmail(this.state.key);

        Promise.resolve(confirmed)
            .then(confirmed => {
                this.setState({
                    loading: false,
                    successful: true,
                });

                let user = localStorage.getItem('user') 
                    ? JSON.parse(localStorage.getItem('user')) 
                    : null;

                if (user === null) {
                    setTimeout(() => {
                        this.props.history('/login')
                    },2000);
                } else {
                    user.is_email_confirmed = "1";
                    localStorage.setItem("user", JSON.stringify(user));

                    setTimeout(() => {
                        this.props.history('/dashboard')
                    },2000);
                }
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
    }

    render() {
        const { isDarkThemeEnabled } = this.state;

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
                        <Col md="9" lg="7" xl="6">
                            <h2 className="alternative text-center mb-3">Confirm Email</h2>
                            <Card className="onboardingCard mb-5">
                                <Card.Body className="p-4">
                                    {
                                        this.state.loading === true
                                        && <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                    }

                                    {
                                        this.state.loading === false
                                        && this.state.successful === false
                                        && "Sorry, something went wrong and email is not confirmed"
                                    }

                                    {
                                        this.state.loading === false
                                        && this.state.successful === true
                                        && "Thank you, your email is confirmed! You will be redirected now."
                                    }
                                </Card.Body>
                            </Card>
                        </Col>
                    </Row>
                </Container>
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

export default connect(mapStateToProps)(withParams(ConfirmEmail));
