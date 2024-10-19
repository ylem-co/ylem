import React, { Component } from 'react';

import Form from 'react-bootstrap/Form';

import PipelineService, {PIPELINE_TYPE_METRIC, PIPELINE_TYPE_GENERIC} from "../services/pipeline.service";
import TaskService from "../services/task.service";

class HeaderSearch extends Component {
    constructor(props) {
        super(props);

        this.wrapperRef = React.createRef();

        this.onChangeString = this.onChangeString.bind(this);
        this.handleSearch = this.handleSearch.bind(this);
        this.handleGetPipelines = this.handleGetPipelines.bind(this);
        this.handleGetTasks = this.handleGetTasks.bind(this);
        this.openSearchResults = this.openSearchResults.bind(this);
        this.hideSearchResults = this.hideSearchResults.bind(this);
        this.handleClickOutside = this.handleClickOutside.bind(this);
        this.handleRedirect = this.handleRedirect.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            string: "",
            pipelines: null,
            metrics: null,
            tasks: null,
            areSearchResultsVisible: false,
        };
    }

    componentDidMount() {
        document.addEventListener("mousedown", this.handleClickOutside);
    }

    componentWillUnmount() {
        document.removeEventListener("mousedown", this.handleClickOutside);
    }

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    handleRedirect(link) {
        window.location.href = link;
    }

    handleClickOutside(event) {
        if (
            this.wrapperRef 
            && this.wrapperRef.current !== null
            && !this.wrapperRef.current.contains(event.target)
        ) {
            this.hideSearchResults();
        }
    }

    onChangeString(e) {
        this.setState({
            string: e.target.value,
            areSearchResultsVisible: true,
        });

        if (e.target.value.length >= 3) {
            this.handleSearch(e.target.value);
        } else if (e.target.value.length === 0) {
            this.hideSearchResults();
        }
    }

    openSearchResults() {
        this.setState({
            areSearchResultsVisible: true,
        });
    }

    hideSearchResults() {
        this.setState({
            areSearchResultsVisible: false,
        });
    }

    handleSearch(string) {
        this.handleGetPipelines(string);
        this.handleGetTasks(string);
    }

    handleGetPipelines = async(string) => {
        var items = PipelineService.searchPipelines(this.state.organization.uuid, string);

        await Promise.resolve(items)
            .then(async(items) => {
                if (items.data && items.data.items && items.data.items !== null) {
                    var pipelines = items.data.items.filter(k => k.type === PIPELINE_TYPE_GENERIC);
                    var metrics = items.data.items.filter(k => k.type === PIPELINE_TYPE_METRIC);

                    await this.promisedSetState({pipelines, metrics});
                } else {
                    await this.promisedSetState({pipelines: [], metrics: []});
                }
            });
    };

    handleGetTasks = async(string) => {
        var tasks = TaskService.searchTasks(this.state.organization.uuid, string);

        await Promise.resolve(tasks)
            .then(async(tasks) => {
                if (tasks.data && tasks.data.items && tasks.data.items !== null) {
                    await this.promisedSetState({tasks: tasks.data.items});
                } else {
                    await this.promisedSetState({tasks: []});
                }
            });
    };

    highlightText(text, highlight) {
        // Split on highlight term and include term into parts, ignore case
        const parts = text.split(new RegExp(`(${highlight})`, 'gi'));
        return <span> { parts.map((part, i) => 
            <span key={i} style={part.toLowerCase() === highlight.toLowerCase() ? { fontWeight: 'bold' } : {} }>
                { part }
            </span>)
        } </span>;
    }

    render() {
        const { string, pipelines, metrics, tasks, areSearchResultsVisible } = this.state;

        return (
            <div ref={this.wrapperRef}>
                <Form.Control
                    className="form-control form-control-sm"
                    type="text"
                    placeholder="&#x1F50D;"
                    id="floatingString"
                    autoComplete="off"
                    name="string"
                    value={string}
                    onChange={this.onChangeString}
                    onClick={this.openSearchResults}
                />
                {
                   (
                        (pipelines !== null && pipelines.length > 0)
                        || (metrics !== null && metrics.length > 0)
                        || (tasks !== null && tasks.length > 0)
                    )
                   && areSearchResultsVisible === true
                   &&
                   <div className="searchResults">
                        {
                            pipelines !== null && pipelines.length > 0
                            &&
                            <div>
                                <div className="searchRow">
                                    <strong>Pipelines</strong>
                                </div>
                                {this.state.pipelines.map(value => (
                                    <div 
                                        className="searchRow" 
                                        key={"w_" + value.uuid}
                                        onClick={(e) => this.handleRedirect(
                                            '/pipelines/folder/' + (value.folder_uuid === "" ? 'root' : value.folder_uuid)
                                            + '/preview/' + value.uuid, e
                                        )}
                                    >
                                        <a href={
                                            '/pipelines/folder/' + (value.folder_uuid === "" ? 'root' : value.folder_uuid)
                                            + '/preview/' + value.uuid
                                        }>
                                            {this.highlightText(value.name, string)}
                                        </a>
                                    </div>
                                ))}
                            </div>
                        }
                        {
                            metrics !== null && metrics.length > 0
                            &&
                            <div>
                                <div className="searchRow">
                                    <strong>Metrics</strong>
                                </div>
                                {this.state.metrics.map(value => (
                                    <div 
                                        className="searchRow" 
                                        key={"m_" + value.uuid}
                                        onClick={(e) => this.handleRedirect(
                                            '/metrics/folder/' + (value.folder_uuid === "" ? 'root' : value.folder_uuid)
                                            + '/details/' + value.uuid, e
                                        )}
                                    >
                                        <a href={
                                            '/metrics/folder/' + (value.folder_uuid === "" ? 'root' : value.folder_uuid)
                                            + '/details/' + value.uuid
                                        }>
                                            {this.highlightText(value.name, string)}
                                        </a>
                                    </div>
                                ))}
                            </div>
                        }
                        {
                            tasks !== null && tasks.length > 0
                            &&
                            <div>
                                <div className="searchRow">
                                    <strong>Tasks</strong>
                                </div>
                                {this.state.tasks.map(value => (
                                    <div 
                                        className="searchRow" 
                                        key={"t_" + value.pipeline_uuid}
                                        onClick={(e) => this.handleRedirect(
                                            '/pipelines/folder/' + (value.folder_uuid === "" ? 'root' : value.folder_uuid)
                                            + '/details/' + value.pipeline_uuid, e
                                        )}
                                    >
                                        <a href={
                                            '/pipelines/folder/' + (value.folder_uuid === "" ? 'root' : value.folder_uuid)
                                            + '/details/' + value.pipeline_uuid
                                        }>
                                            {this.highlightText(value.name, string)}
                                        </a>
                                    </div>
                                ))}
                            </div>
                        }
                   </div>
                }
            </div>
        );
    }
}

export default HeaderSearch;
