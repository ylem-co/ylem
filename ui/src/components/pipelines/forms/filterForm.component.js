import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";

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
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap';

export const equations = {
  neq:  '!=',
  eq:   '==',
  let:  '<=',
  het:  '>=',
  lt:   '<',
  ht:   '>',
};

class FilterForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeFilterField = this.onChangeFilterField.bind(this);
        this.onChangeFilterOperation = this.onChangeFilterOperation.bind(this);
        this.onChangeFilterValue = this.onChangeFilterValue.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: this.props.item.name,
            severity: TASK_SEVERITY_MEDIUM,
            expression: "",
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
            filterOperation: "eq",
            filterValue: "",
            filterField: "",
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

    onChangeFilterValue = async(e) => {
        await this.promisedSetState({
            filterValue: e.target.value,
        });
    }

    onChangeFilterField = async(e) => {
        await this.promisedSetState({
            filterField: e.target.value,
        });
    }

    onChangeFilterOperation = async(key) => {
        await this.promisedSetState({
            filterOperation: key,
        });
    }

    mapFormToItem() {
        let value = this.state.filterValue;

        if (isNaN(value) || value === "") {
            value = '"' + value + '"';
        }

        let data = {
            "type": this.state.type,
            "expression": '#(' + this.state.filterField + ' ' + equations[this.state.filterOperation] + ' ' + value + ')#',
        };

        return data;
    };

    mapItemToForm = async(item) => {
        let filterField = "";
        let filterOperation = "eq";
        let filterValue = "";

        const re = new RegExp(/\#\((.+)\s(.+)\s\"*([^\"]*)\"*\)\#/i);
        let matches = re.exec(item);

        if (matches !== null) {
            filterField = matches[1];
            filterOperation = Object.keys(equations).find(key => equations[key] === matches[2]);
            if (matches[3] !== null) {
                filterValue = matches[3];
            }
        }

        await this.promisedSetState({filterField, filterOperation, filterValue});
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
            dispatch(
                updateTask(
                    this.state.item.uuid, 
                    this.state.item.pipeline_uuid,
                    this.state.name,
                    this.state.severity,
                    this.state.item.type,
                    this.mapFormToItem()
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
                                <FloatingLabel controlId="floatingField block" label="Field">
                                    <Input
                                        className="form-control form-control-lg"
                                        type="text"
                                        id="floatingField"
                                        placeholder="Field"
                                        autoComplete="field"
                                        name="field"
                                        value={this.state.filterField}
                                        onChange={this.onChangeFilterField}
                                        validations={[required]}
                                    />
                                </FloatingLabel>
                                <Dropdown size="lg" className="mb-4 block">
                                    <Dropdown.Toggle 
                                        variant="light" 
                                        id="dropdown-basic"
                                        className="mt-2 block"
                                    >
                                        {
                                            this.state.filterOperation !== null
                                                ? equations[this.state.filterOperation]
                                                : equations.neq
                                        }
                                    </Dropdown.Toggle>

                                    <Dropdown.Menu>
                                        {Object.keys(equations).map((key, index) => (
                                            <Dropdown.Item
                                                value={equations[key]}
                                                key={"col" + index}
                                                active={equations[key] === this.state.filterOperation}
                                                onClick={(e) => this.onChangeFilterOperation(key)}
                                            >
                                                {equations[key]}
                                            </Dropdown.Item>
                                        ))}
                                    </Dropdown.Menu>
                                </Dropdown>
                                <FloatingLabel controlId="floatingValue block" label="Value">
                                    <Input
                                        className="form-control form-control-lg"
                                        type="text"
                                        id="floatingValue"
                                        placeholder="Value"
                                        autoComplete="value"
                                        name="value"
                                        value={this.state.filterValue}
                                        onChange={this.onChangeFilterValue}
                                    />
                                </FloatingLabel>
                            <div className="inputTip">For more advanced filtering options, use the "Transformer" task with the "Extract from JSON" option.</div>
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

export default connect(mapStateToProps)(FilterForm);
