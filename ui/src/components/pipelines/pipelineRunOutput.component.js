import React from 'react';

import Spinner from "react-bootstrap/Spinner";

class PipelineRunOutput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {
        const { finished, output } = this.props;

        return (
                <div className="CLIImitation">
                    {output}
                    {
                        finished !== true
                            && <div className="text-center mt-3"><Spinner animation="grow" className="spinner-secondary"/></div> 
                    }
                </div>
        );
    }
}

export default PipelineRunOutput;
