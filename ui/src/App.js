import React, { Component } from "react";
import { connect } from "react-redux";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import Login from "./components/login.component";
import Register from "./components/register.component";
import ConfirmEmail from "./components/confirmEmail.component";
import Invitation from "./components/invitation.component";

import axios from "axios";
import { v4 as uuidv4 } from 'uuid';

import { logout } from "./actions/auth";
import { history } from './helpers/history';
import { clearMessage } from "./actions/message";

import './App.scss';
import './Modal.scss';
import ExternalSignIn from "./components/externalSignIn.component";
//import 'bootstrap/dist/css/bootstrap.min.css';

// Containers
const Layout = React.lazy(() => import('./containers/layout'));

const loading = (
  <div className="pt-3 text-center">
    <div className="sk-spinner sk-spinner-pulse"></div>
  </div>
)

/*if (process.env.REACT_APP_BACKEND_URL) {
  axios.defaults.baseURL = window.location.protocol + process.env.REACT_APP_BACKEND_URL;
}*/

class App extends Component {
  constructor(props) {
    super(props);
    this.logOut = this.logOut.bind(this);

    this.state = {
      currentUser: undefined,
      isTourOpen: window.innerWidth > 992 && !localStorage.getItem('tourIsWatched'),
    };

    history.listen((location) => {
      props.dispatch(clearMessage()); // clear message when changing location
    });
  }

  componentDidMount() {
    const user = this.props.user;

    if (user) {
      this.setState({
        currentUser: user,
      });
    }

    // Set default API call headers
    axios.interceptors.request.use(
      config => {
        config.headers['X-Request-Id'] = uuidv4();
        var token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = 'Bearer ' + token;
        }

        return config;
      }
    );

    // impersonation request interceptor
    axios.interceptors.request.use((config) => {
      var impersonating = localStorage.getItem('impersonating') || null;

      if (impersonating !== null) {
        config.params = config.params || {};
        config.params['_switch_user'] = impersonating;
      }
      return config;
    });

    // token expired or invalid token handler
    axios.interceptors.response.use(response => {
      return response;
    }, error => {
      if (
          error.response && error.response.status === 401
          //&& error.response.data.message === 'Expired JWT Token'
      ) {
        if (error.response.data.message === "Authentication request could not be processed due to a system problem.") {
          this.emailIsNotYetConfirmed();
        } else if (error.response.data.message === "Authorization Failed") {
          //this.logOut();
        } else {
          //this.logOut();
        }
      }

      if (
          error.response && error.response.status === 477
      ) {
        this.emailIsNotYetConfirmed();
      }

      return Promise.reject(error);
    });
  }

  logOut() {
    this.props.dispatch(logout());
  }

  emailIsNotYetConfirmed() {
    history.push("/477");
  }

  closeTour = () => {
    this.setState({isTourOpen: false});
    localStorage.setItem('tourIsWatched', true);
  };

  render() {

    return (
        <>
          <Router history={history} basename="/">
            <React.Suspense fallback={loading}>
              <Routes>
                <Route path="/confirm-email/:key" name="Confirm Email" element={<ConfirmEmail/>} />
                <Route path="/invitation/:key" name="Invitation by Key" element={<Invitation/>} />
                <Route path="/login" name="Login Page" element={<Login/>} />
                <Route path="/register" name="Register Page" element={<Register/>} />
                <Route path="/sign-in/google" name="Sign-In via Google" element={<ExternalSignIn/>} />
                <Route path="*" name="Home" element={<Layout/>} />
              </Routes>
            </React.Suspense>
          </Router>
        </>
    );
  }
}

function mapStateToProps(state) {
  const { user } = state.auth;
  return {
    user,
  };
}

export default connect(mapStateToProps)(App)
