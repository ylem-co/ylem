import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";

import CodeEditor from '@uiw/react-textarea-code-editor';
import rehypePrism from 'rehype-prism-plus';

import { connect } from "react-redux";
import { updateTask } from "../../../actions/tasks";

import Input from "../../formControls/input.component";
import { required } from "../../formControls/validations";

import { clearMessage } from "../../../actions/message";

import { 
    TASK_SEVERITY_MEDIUM,
    TASK_PROCESSOR_STRATEGY_INCLUSIVE,
} from "../../../services/task.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap'

class ProcessorForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeStrategy = this.onChangeStrategy.bind(this);
        this.onChangeExpression = this.onChangeExpression.bind(this);
        this.onChangeExpressionFromOutside = this.onChangeExpressionFromOutside.bind(this);
        this.onChangeSeverity = this.onChangeSeverity.bind(this);

        this.state = {
            name: "",
            strategy: TASK_PROCESSOR_STRATEGY_INCLUSIVE,
            severity: TASK_SEVERITY_MEDIUM,
            expression: "",
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
                expression: this.props.item.implementation.expression || "",
                strategy: this.props.item.implementation.strategy || TASK_PROCESSOR_STRATEGY_INCLUSIVE,
            });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeStrategy(strategy) {
        this.setState({strategy});
    }

    onChangeExpression(e) {
        this.setState({
            expression: e.target.value,
        });
    }

    onChangeExpressionFromOutside(el) {
        this.setState({
            expression: el.value,
        });
    }

    onChangeSeverity(severity) {
        this.setState({ severity })
    }

    handleSubmit(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            let name = this.state.name;
            if (name === "") {
                name = this.state.expression;
            }

            let implementation = {
                "expression": this.state.expression,
            }

            if (this.state.strategy !== "") {
                implementation.strategy = this.state.strategy;
            }

            dispatch(
                updateTask(
                    this.state.item.uuid, 
                    this.state.item.pipeline_uuid,
                    name,
                    this.state.severity,
                    this.state.item.type,
                    implementation
                )
            )
            .then(() => {
                var item = this.state.item;
                item.name = name;
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
                            <FloatingLabel controlId="floatingName" label="Name (optional)">
                                <Input
                                    className="form-control form-control-lg"
                                    type="text"
                                    id="floatingName"
                                    placeholder="Name (optional)"
                                    autoComplete="name"
                                    name="name"
                                    value={this.state.name}
                                    onChange={this.onChangeName}
                                />
                            </FloatingLabel>
                        </div>
                    </InputGroup>

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <span className="inputLabel">Processor</span><br/>
                            <Dropdown size="lg" className="mb-4">
                                <Dropdown.Toggle
                                    variant="light"
                                    disabled={true}
                                    id="dropdown-basic"
                                >
                                    JQ
                                </Dropdown.Toggle>
                            </Dropdown>
                        </div>
                    </InputGroup>

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <label className="nonFloatingLabel">Expression</label>
                                <CodeEditor
                                    className="form-control form-control-lg codeEditor"
                                    type="text"
                                    language="xls"
                                    id="floatingExpression"
                                    minHeight={200}
                                    autoComplete="Expression"
                                    name="expression"
                                    value={this.state.expression}
                                    onChange={this.onChangeExpression}
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
                                    For processing, we support the functionality of the lightweight JSON processor <a href="https://jqlang.github.io/jq/" target="_blank" rel="noreferrer">JQ</a>.<br/><br/>
                                    Example:<br/>
                                    {'{'}<br/>
                                    &nbsp;&nbsp;&nbsp;&nbsp;IDs: map(.id),<br/>
                                    &nbsp;&nbsp;&nbsp;&nbsp;totalAmount: map(.amount) | add,<br/>
                                    &nbsp;&nbsp;&nbsp;&nbsp;totalAmountInEUR: map(select(.currency == "EUR").amount) | add<br/>
                                    {'}'}<br/><br/>
                                    Check their <a href="https://jqlang.github.io/jq/manual/" target="_blank" rel="noreferrer">documentation</a> for more syntax.
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

export default connect(mapStateToProps)(ProcessorForm);
