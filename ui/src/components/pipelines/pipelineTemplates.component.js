import React from 'react';
import { useNavigate } from 'react-router-dom';

import Spinner from "react-bootstrap/Spinner";
import Nav from 'react-bootstrap/Nav';

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import FloatingLabel from "react-bootstrap/FloatingLabel";
import Button from "react-bootstrap/Button";

import PipelineService, {PIPELINE_TYPE_METRIC} from "../../services/pipeline.service";
import TemplateService from "../../services/template.service";

import CorporateFareRounded from '@mui/icons-material/CorporateFareRounded';
import PersonOutlineRounded from '@mui/icons-material/PersonOutlineRounded';
import SettingsOutlined from '@mui/icons-material/SettingsOutlined';

import {
    TEMPLATES_LIST_TYPE_SYSTEM,
    TEMPLATES_LIST_TYPE_ORG,
    TEMPLATES_LIST_TYPE_ME,
} from "../../services/template.service";

import {
    PIPELINE_PAGE_DETAILS
} from "../../services/pipeline.service";

function withParams(Component) {
  return props => <Component {...props} history={useNavigate()} />;
}

class PipelineTemplates extends React.Component {
    constructor(props) {
        super(props);
        this.onChangeActiveTab = this.onChangeActiveTab.bind(this);
        this.handleGetTemplates = this.handleGetTemplates.bind(this);
        this.handleGetPreviews = this.handleGetPreviews.bind(this);
        this.activateTemplateUuid = this.activateTemplateUuid.bind(this);
        this.createFromTemplate = this.createFromTemplate.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : null,
            elements: null,
            activeTab: TEMPLATES_LIST_TYPE_SYSTEM,
            activeTemplateUuid: null,
            loading: false,
            errorMessage: null,
        };
    }

    componentDidMount() {
        let elements = this.handleGetTemplates("system");
        this.setState({elements});
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    createFromTemplate = async() => {
        await this.promisedSetState({loading: true});
        let pipeline = TemplateService.createFromTemplate(
            this.state.activeTemplateUuid,
            this.props.folderUuid
        );

        await Promise.resolve(pipeline)
            .then(async(pipeline) => {
                if (pipeline.uuid) {
                    this.props.history(
                        this.props.type === PIPELINE_TYPE_METRIC
                            ? '/metrics/folder/' + this.props.folderUuid + '/' + PIPELINE_PAGE_DETAILS + '/' + pipeline.uuid
                            : '/pipelines/folder/' + this.props.folderUuid + '/' + PIPELINE_PAGE_DETAILS + '/' + pipeline.uuid
                    );

                    window.location.reload();
                }
            })
            .catch((error) => {
                this.setState({
                    errorMessage: error.response.data.message,
                    successful: false,
                    loading: false,
                });
            });
    }

    activateTemplateUuid(activeTemplateUuid) {
        this.setState({
            activeTemplateUuid: activeTemplateUuid !== this.state.activeTemplateUuid
                ? activeTemplateUuid
                : null
        });
    };

    onChangeActiveTab = async(activeTab) => {
        await this.promisedSetState({
            activeTab,
            elements: null,
        });

        let elements = this.handleGetTemplates(activeTab);
        await this.promisedSetState({elements});
    };

    handleGetTemplates = async(type) => {
        let templates = TemplateService.getTemplates(
            type, 
            this.state.organization !== null
                ? this.state.organization.uuid
                : null
        );

        await Promise.resolve(templates)
            .then(async(templates) => {
                if (templates.data && templates.data.items !== null) {
                    var items = templates.data.items;

                    items = items.filter(k => k.type === this.props.type);

                    await items.sort((a, b) => a.name.localeCompare(b.name));
   
                    await this.promisedSetState({elements: items});

                    if (this.props.type !== PIPELINE_TYPE_METRIC) {
                        this.handleGetPreviews(items);
                    }
                } else {
                    await this.promisedSetState({elements: []});
                }
            })
            .catch(async() => {
                await this.promisedSetState({elements: []});
            });
    };

    handleGetPreviews = async(items) => {
        var image = "";
        var src = "";

        for(var i = 0; i < items.length; i++){
            if (items[i].preview !== "") {
                image = PipelineService.getPipelinePreview(items[i].uuid, true);
                src = await Promise.resolve(image)
                .then(value => {
                    return value;
                });
                items[i].src = src;
                this.setState({elements: items});
            }
        }
    };

    search = (e) => {
        let searchString = e.target.value.toLowerCase();
        const container = document.getElementById("templatesListContainer").firstChild;

        for(var child = container.firstChild; child !== null; child = child.nextSibling) {
            child.classList.remove("hiddenTemplateInTheList");
            if (child.firstChild.firstChild.firstChild.firstChild.innerText.toLowerCase().indexOf(searchString) === -1) {
                child.classList.add("hiddenTemplateInTheList");
            }
        }
    };

    render() {
        const { activeTab, errorMessage, loading, successful } = this.state;

        return (
            <div>
                <Nav variant="tabs" className="mb-4">
                    <Nav.Item>
                        <Nav.Link
                            active={activeTab === TEMPLATES_LIST_TYPE_SYSTEM}
                            onClick={() => this.onChangeActiveTab(TEMPLATES_LIST_TYPE_SYSTEM)}
                            className="tmpList"
                        >
                            <SettingsOutlined className="tabIcon"/> System templates
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link
                            active={activeTab === TEMPLATES_LIST_TYPE_ORG}
                            onClick={() => this.onChangeActiveTab(TEMPLATES_LIST_TYPE_ORG)}
                            className="tmpList"
                        >
                            <CorporateFareRounded className="tabIcon"/> My organization's templates
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link
                            active={activeTab === TEMPLATES_LIST_TYPE_ME}
                            onClick={() => this.onChangeActiveTab(TEMPLATES_LIST_TYPE_ME)}
                            className="tmpList"
                        >
                            <PersonOutlineRounded className="tabIcon"/> My templates
                        </Nav.Link>
                    </Nav.Item>
                </Nav>
                <div className="mb-4">
                    <FloatingLabel controlId="floatingSearch" label="Search templates">
                        <input
                            onChange={(e) => this.search(e)}
                            type="text"
                            id="floatingSearch"
                            className="form-control form-control-lg"
                        />
                    </FloatingLabel>
                </div>
                <div id="templatesListContainer">
                {
                    this.state.elements !== null ?
                        this.state.elements.length > 0 ?
                            <Row>
                                {this.state.elements.map(value => (
                                    value.is_template === 1 &&
                                        <Col className="col-3 mb-4" key={value.uuid}>
                                            <Card 
                                                className={
                                                    value.uuid !== this.state.activeTemplateUuid
                                                    ? "withHeader onHoverCard"
                                                    : "withHeader onHoverCard activeTemplate"
                                                }
                                                onClick={() => this.activateTemplateUuid(value.uuid)}
                                            >
                                                <Card.Header className="noTopBorder">
                                                    <Row>
                                                        <Col className="col-12">
                                                            {value.name}
                                                        </Col>
                                                    </Row>
                                                </Card.Header> 
                                                <Card.Body>
                                                    { this.props.type === PIPELINE_TYPE_METRIC
                                                        ?
                                                            <div className="mb-3"></div>
                                                        :
                                                            <div className="mb-3">
                                                                {
                                                                    'src' in value ?
                                                                        value.src !== ""
                                                                        ? <div className="pipelineInTheList"><img src={'data:image/jpeg;base64,' + value.src} className="pipelinePreview" alt={value.name} title={value.name}/></div>
                                                                        : <div></div>
                                                                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                                                }
                                                            </div>
                                                    }
                                                </Card.Body>
                                            </Card>
                                        </Col>
                                    ))} 
                            </Row>
                        : <div>No templates found in this category</div>
                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                </div>
                {errorMessage !== null 
                    && loading === false 
                    && successful === false 
                    && (
                  <div className="form-group">
                    <div className={ successful ? "alert alert-success" : "alert alert-danger" } role="alert">
                      {errorMessage === null ? "Pipeline successfully created" : errorMessage}
                    </div>
                  </div>
                )}
                <div className="text-right">
                    <Button
                        variant="secondary"
                        className="mx-0 mt-3"
                        disabled={this.state.loading || this.state.activeTemplateUuid === null}
                        onClick={() => this.createFromTemplate()}
                    >
                        {this.state.loading && (
                            <span className="spinner-border spinner-border-sm spinner-primary"></span>
                        )}
                        <span>Create</span>
                    </Button>
                </div>
            </div>
        );
    }
}

export default withParams(PipelineTemplates);
