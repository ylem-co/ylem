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
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap'

const TRANSFORMER_TYPE_STR_SPLIT = "str_split";
const TRANSFORMER_TYPE_EXTRACT_FROM_JSON = "extract_from_json";
const TRANSFORMER_TYPE_CAST_TO = "cast_to";
const TRANSFORMER_TYPE_ENCODE_TO = "encode_to";

const TransformerTypes = {
    str_split: "Split string",
    extract_from_json: "Extract from JSON",
    cast_to: "Cast to",
    encode_to: "Encode to",
};

const ENCODE_FORMAT_XML = "XML";
const ENCODE_FORMAT_CSV = "CSV";

const EncodeFormats = [
    ENCODE_FORMAT_XML,
    ENCODE_FORMAT_CSV,
];

const CAST_TO_TYPE_STRING = "string";

const CastToTypes = {
    string: "String",
    integer: "Integer",
};

class TransformerForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeType = this.onChangeType.bind(this);
        this.onChangeJsonQueryExpression = this.onChangeJsonQueryExpression.bind(this);
        this.onChangeDelimiter = this.onChangeDelimiter.bind(this);
        this.onChangeCastToType = this.onChangeCastToType.bind(this);
        this.onChangeDecodeFormat = this.onChangeDecodeFormat.bind(this);
        this.onChangeEncodeFormat = this.onChangeEncodeFormat.bind(this);
        this.onChangeSeverity = this.onChangeSeverity.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: this.props.item.name,
            severity: TASK_SEVERITY_MEDIUM,
            jsonQueryExpression: "",
            delimiter: "",
            castToType: CAST_TO_TYPE_STRING,
            castToTypeValue: CastToTypes.string,
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
            type: TRANSFORMER_TYPE_STR_SPLIT,
            typeValue: TransformerTypes.str_split,
            decodeFormat: "",
            encodeFormat: ENCODE_FORMAT_CSV,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
                jsonQueryExpression: this.props.item.implementation.json_query_expression || "",
                type: this.props.item.implementation.type || TRANSFORMER_TYPE_STR_SPLIT,
                typeValue: this.props.item.implementation.type ? TransformerTypes[this.props.item.implementation.type] : TransformerTypes.str_split,
                delimiter: this.props.item.implementation.delimiter || "",
                castToType: this.props.item.implementation.cast_to_type || CAST_TO_TYPE_STRING,
                castToTypeValue: this.props.item.implementation.cast_to_type ? CastToTypes[this.props.item.implementation.cast_to_type] : CastToTypes.string,
                encodeFormat: this.props.item.implementation.encode_format || ENCODE_FORMAT_CSV,
                decodeFormat: this.props.item.implementation.decode_format || "",
            });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeJsonQueryExpression(e) {
        this.setState({
            jsonQueryExpression: e.target.value,
        });
    }

    onChangeDelimiter(e) {
        this.setState({
            delimiter: e.target.value,
        });
    }

    onChangeCastToType(type) {
        this.setState({
            castToType: type,
            castToTypeValue: CastToTypes[type],
        });
    }

    onChangeEncodeFormat(format) {
        this.setState({
            encodeFormat: format,
        });
    }

    onChangeDecodeFormat(format) {
        this.setState({
            decodeFormat: format,
        });
    }

    onChangeType(type) {
        this.setState({
            type: type,
            typeValue: TransformerTypes[type],
        })
    }

    onChangeSeverity(severity) {
        this.setState({ severity })
    }

    mapFormToItem() {
        let data = {
            "type": this.state.type,
            "json_query_expression": "",
            "delimiter": "",
            "cast_to_type": "",
            "decode_format": "",
            "encode_format": "",
        };

        if (this.state.type === TRANSFORMER_TYPE_STR_SPLIT) {
            data.delimiter = this.state.delimiter;
        }

        if (this.state.type === TRANSFORMER_TYPE_CAST_TO) {
            data.cast_to_type = this.state.castToType;
        }

        if (this.state.type === TRANSFORMER_TYPE_EXTRACT_FROM_JSON) {
            data.json_query_expression = this.state.jsonQueryExpression;
        }

        if (this.state.type === TRANSFORMER_TYPE_ENCODE_TO) {
            if (this.state.encodeFormat === ENCODE_FORMAT_CSV) {
                data.delimiter = this.state.delimiter;
            }
            data.encode_format = this.state.encodeFormat;
        }

        return data;
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

                    <span className="inputLabel">Type</span><br/>
                    <Dropdown size="lg" className="mb-4">
                        <Dropdown.Toggle 
                            variant="light" 
                            id="dropdown-basic"
                        >
                            {this.state.typeValue}
                        </Dropdown.Toggle>

                        <Dropdown.Menu>
                            {Object.entries(TransformerTypes).map(([key, value]) => (
                                <Dropdown.Item
                                    value={key}
                                    key={key}
                                    active={key === this.state.type}
                                    onClick={(e) => this.onChangeType(key)}
                                >
                                    {value}
                                </Dropdown.Item>
                            ))}
                        </Dropdown.Menu>
                    </Dropdown>

                    { this.state.type === TRANSFORMER_TYPE_ENCODE_TO &&
                        <>
                    <span className="inputLabel">Encoding format</span><br/>
                    <Dropdown size="lg" className="mb-4">
                        <Dropdown.Toggle 
                            variant="light" 
                            id="dropdown-basic"
                        >
                            {this.state.encodeFormat}
                        </Dropdown.Toggle>

                        <Dropdown.Menu>
                            {Object.entries(EncodeFormats).map(([key, value]) => (
                                <Dropdown.Item
                                    value={value}
                                    key={value}
                                    active={value === this.state.encodeFormat}
                                    onClick={(e) => this.onChangeEncodeFormat(value)}
                                >
                                    {value}
                                </Dropdown.Item>
                            ))}
                        </Dropdown.Menu>
                    </Dropdown>
                    </>
                }

                    { 
                        (
                            (
                                this.state.encodeFormat === ENCODE_FORMAT_CSV 
                                && this.state.type === TRANSFORMER_TYPE_ENCODE_TO
                            )
                        || this.state.type === TRANSFORMER_TYPE_STR_SPLIT) &&
                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <FloatingLabel controlId="floatingDelimiter" label="Delimiter">
                                <Input
                                    className="form-control form-control-lg"
                                    type="text"
                                    id="floatingDelimiter"
                                    placeholder="Delimiter"
                                    autoComplete="delimiter"
                                    name="delimiter"
                                    value={this.state.delimiter}
                                    onChange={this.onChangeDelimiter}
                                />
                                <div className="inputTip">In case you decode or encode .csv or perform a string split, specify a delimiter</div>
                            </FloatingLabel>
                        </div>
                    </InputGroup>
                    }

                    { this.state.type === TRANSFORMER_TYPE_CAST_TO &&
                        <>
                            <span className="inputLabel">Type</span><br/>
                            <Dropdown size="lg" className="mb-4">
                                <Dropdown.Toggle 
                                    variant="light" 
                                    id="dropdown-basic"
                                >
                                    {this.state.castToTypeValue}
                                </Dropdown.Toggle>

                                <Dropdown.Menu>
                                    {Object.entries(CastToTypes).map(([key, value]) => (
                                        <Dropdown.Item
                                            value={key}
                                            key={key}
                                            active={key === this.state.castToType}
                                            onClick={(e) => this.onChangeCastToType(key)}
                                        >
                                            {value}
                                        </Dropdown.Item>
                                    ))}
                                </Dropdown.Menu>
                            </Dropdown>
                        </>
                    }

                    { 
                        this.state.type === TRANSFORMER_TYPE_EXTRACT_FROM_JSON
                        &&
                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <span className="inputLabel">JSON Query</span>
                            <CodeEditor
                                className="form-control form-control-lg codeEditor"
                                type="text"
                                id="floatingJsonQueryExpression"
                                autoComplete="jsonQueryExpression"
                                name="jsonQueryExpression"
                                value={this.state.jsonQueryExpression}
                                onChange={this.onChangeJsonQueryExpression}
                                language="powerquery"
                                minHeight={200}
                                rehypePlugins={[
                                    [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                ]}
                                style={{
                                    fontSize: 14,
                                    fontFamily: 'Source Code Pro, monospace',
                                }}
                            />
                            <div className="inputTip">We support JSON queries based on the functionality of <a href="https://github.com/tidwall/gjson" target="_blank" rel="noreferrer">GJSON</a> library</div>
                        </div>
                    </InputGroup>
                    }

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

export default connect(mapStateToProps)(TransformerForm);
