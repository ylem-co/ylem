import React, { Component } from 'react';

import Nav from 'react-bootstrap/Nav';

import { 
    PIPELINE_TYPE_METRIC,
    PIPELINE_PAGE_STATS,
    PIPELINE_PAGE_LOGS,
    PIPELINE_PAGE_DETAILS,
    PIPELINE_PAGE_TRIGGERS,
    PIPELINE_PAGE_PREVIEW,
} from "../../services/pipeline.service";

class PipelineTabs extends Component {
    constructor(props) {
        super(props);

        this.state = {
            page: this.props.page,
            type: this.props.type === PIPELINE_TYPE_METRIC ? PIPELINE_TYPE_METRIC : "preview",
            statLink: this.props.statLink,
            pipelineLink: this.props.pipelineLink,
            triggerLink: this.props.triggerLink,
            logsLink: this.props.logsLink,
        };
    }

    render() {
        return (
            <>
                <Nav variant="tabs" className="mb-4">
                    <Nav.Item>
                        <Nav.Link
                            href={this.state.pipelineLink}
                            active={
                                this.state.page === PIPELINE_PAGE_DETAILS
                                || this.state.page === PIPELINE_PAGE_PREVIEW
                            }
                            className="pipelineTab"
                        >
                            { 
                                this.state.type.charAt(0).toUpperCase() 
                                + this.state.type.slice(1)
                            }
                        </Nav.Link>
                    </Nav.Item>
                    { this.state.type !== PIPELINE_TYPE_METRIC
                        &&
                        <Nav.Item>
                            <Nav.Link
                                href={this.state.triggerLink}
                                active={this.state.page === PIPELINE_PAGE_TRIGGERS}
                                className="pipelineTab"
                            >
                                Triggers
                            </Nav.Link>
                        </Nav.Item>
                    }
                    <Nav.Item>
                        <Nav.Link
                            href={this.state.statLink}
                            active={this.state.page === PIPELINE_PAGE_STATS}
                            className="pipelineTab"
                        >
                            Statistics
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link
                            href={this.state.logsLink}
                            active={this.state.page === PIPELINE_PAGE_LOGS}
                            className="pipelineTab"
                        >
                            Logs
                        </Nav.Link>
                    </Nav.Item>
                </Nav>
            </>
        );
    }
}

export default PipelineTabs;
