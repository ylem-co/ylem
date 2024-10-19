import React, {Component} from "react";
import {Link, Navigate, useNavigate} from 'react-router-dom';
import { Fade } from "react-awesome-reveal";
import {connect} from "react-redux";
import {signInWithGoogle} from "../actions/auth";
import {PERMISSION_LOGGED_OUT, validatePermissions} from "../actions/pipeline";
import {Card, CardGroup, Col, Container, Row} from 'react-bootstrap'

function withParams(Component) {
  return props => <Component {...props} history={useNavigate()} />;
}

class ExternalSignIn extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loading: true,
            isDarkThemeEnabled: localStorage.getItem('darkTheme') !== "false",
        };
    }

    componentDidMount() {
        document.title = 'Signing you in...'

        const { dispatch, history } = this.props;

        dispatch(signInWithGoogle(window.location.search))
            .then(() => {
                history("/dashboard");
            })
            .catch(() => {
                this.setState({
                    loading: false,
                });
            });
    };

    render() {
        const { isLoggedIn, user, message } = this.props;

        const { isDarkThemeEnabled } = this.state;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/dashboard" />
        }

        return (
        <div className="align-items-center visualBg">
            {this.state.loading && (
                <span className="spinner-border spinner-border-sm spinner-primary spinner-centered spinner-border-lg"></span>
            )}
            {message && (
                <Fade>
                    <Container>
                        <Row className="justify-content-center">
                            <Col md="8" sm="12" className="text-center pt-5">
                                <Link to="/" className="">
                                    {isDarkThemeEnabled
                                        ? <img src="/images/logo2-dark.png" width="150px" alt="Ylem"/>
                                        : <img src="/images/logo2.png" width="150px" alt="Ylem"/>
                                    }
                                </Link>
                            </Col>
                        </Row>
                        <Row className="justify-content-center pt-5">
                            <Col lg="10" md="12" xl="8">
                                <CardGroup>
                                    <Card className="p-4 onboardingCard">
                                        <Card.Body>
                                            <div className="form-group">
                                                <div className="alert alert-danger mt-3" role="alert">
                                                    {message}
                                                </div>
                                            </div>
                                        </Card.Body>
                                    </Card>
                                </CardGroup>
                            </Col>
                        </Row>
                    </Container>
                </Fade>
            )}
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

export default connect(mapStateToProps)(withParams(ExternalSignIn));
