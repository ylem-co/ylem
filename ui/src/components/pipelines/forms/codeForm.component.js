import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";

import CodeEditor from '@uiw/react-textarea-code-editor';
import rehypePrism from 'rehype-prism-plus';

import { connect } from "react-redux";
import { updateTask } from "../../../actions/tasks";

import { clearMessage } from "../../../actions/message";

import Input from "../../formControls/input.component";
import { required } from "../../formControls/validations";

import { TASK_SEVERITY_MEDIUM } from "../../../services/task.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    InputGroup,
    Row
} from 'react-bootstrap';

class CodeForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeCode = this.onChangeCode.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: this.props.item.name,
            severity: TASK_SEVERITY_MEDIUM,
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
            code: "",
        };
    }

    componentDidMount = async() => {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            if (this.props.item.implementation.expression) {
                await this.mapItemToForm(this.props.item.implementation.expression);
            }

            await this.promisedSetState({
                item: this.props.item,
                name: this.props.item.name || "",
                code: this.props.item.implementation.code || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
            });
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeCode(e) {
        this.setState({
            code: e.target.value,
        });
    }

    mapFormToItem() {
        return {
            "type": this.state.type,
        };
    };

    handleSubmit(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                updateTask(
                    this.state.item.uuid, 
                    this.state.item.pipeline_uuid,
                    this.state.name,
                    this.state.severity,
                    this.state.item.type,
                    {
                        "code": this.state.code,
                        "type": "python", // tmp until we have other languages
                    }
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

                this.props.successHandler(item);
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
                    onSubmit={this.handleSubmit}
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

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <label className="nonFloatingLabel">Python code</label>
                            <CodeEditor
                                className="form-control form-control-lg codeEditor"
                                type="text"
                                language="python"
                                minHeight={230}
                                autoComplete="SQLQuery"
                                name="CodeCode"
                                value={this.state.code}
                                onChange={this.onChangeCode}
                                validations={[required]}
                                rehypePlugins={[
                                    [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                ]}
                                style={{
                                    fontSize: 14,
                                    fontFamily: 'Source Code Pro, monospace',
                                }}
                            />
                            <div className="inputTip">
                                Put your python code here. JSON object is available as a variable "input"
                                <br/><br/>
                                # Just an example<br/>
                                {'result = []'}<br/><br/>
                                {'for i in input:'}<br/>
                                &nbsp;&nbsp;{'if i["amount"] is not None and i["amount"] < 2000:'}<br/>
                                &nbsp;&nbsp;&nbsp;&nbsp;{'result.append(i)'}<br/><br/>
                                {'input = result'}
                            </div>
                            <div className="inputTip">
                                You can also use Pandas framework for the operations like this:
                                <br/><br/>
                                df = pandas.DataFrame(input) # that's how you need to create a Pandas data frame from the task input<br/>
                                result = df.head(2) # a basic operation example performed by Pandas<br/>
                                input = result.to_dict() # that's how you create an output of the task from the result data frame
                            </div>
                        </div>
                    </InputGroup>

                    <Row>
                        <Col xs="6">
                            <Button
                                className="px-4 btn btn-primary"
                                disabled={this.state.loading}
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

export default connect(mapStateToProps)(CodeForm);
