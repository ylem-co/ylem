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

class GptForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangePrompt = this.onChangePrompt.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: this.props.item.name,
            severity: TASK_SEVERITY_MEDIUM,
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
            prompt: "",
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
                prompt: this.props.item.implementation.prompt || "",
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

    onChangePrompt(e) {
        this.setState({
            prompt: e.target.value,
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
                        "prompt": this.state.prompt,
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
                            <label className="nonFloatingLabel">Your request to GPT</label>
                                <CodeEditor
                                    className="form-control form-control-lg codeEditor"
                                    type="text"
                                    language="xls"
                                    id="floatingBody"
                                    minHeight={200}
                                    autoComplete="prompt"
                                    name="prompt"
                                    value={this.state.prompt}
                                    onChange={this.onChangePrompt}
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
                                Write what you would like to do to your JSON data in a natural language.<br/><br/>
                                For example:<br/>
                                <ul>
                                    <li>Multiply all amounts by 3 and return new JSON only. Without other text</li>
                                    <li>Extract postal codes from addresses and return new JSON</li>
                                </ul>
                                Due to the limits of ChatGPT, the max length of data is 1024 tokens. If your data is more extended please consider using only part of it.
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

export default connect(mapStateToProps)(GptForm);
