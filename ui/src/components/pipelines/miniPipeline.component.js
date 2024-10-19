import React, { Component } from 'react';

import ReactFlow, {
  Background,
} from 'react-flow-renderer';

import { connect } from "react-redux";

import { nodeTypes } from './pipeline.component';

class MiniPipeline extends Component {
    constructor(props) {
      super(props);

      this.onLoad = this.onLoad.bind(this);

      this.state = {
        elements: null,
        item: this.props.item,
        reactFlowInstance: null,
      };
    }

    componentDidMount() {
        this.setState({
            elements: JSON.parse(this.props.elements),
            item: this.props.item,
        });
    };

    UNSAFE_componentWillReceiveProps(nextProps, nextContext) {
        this.setState({
            elements: JSON.parse(this.props.elements),
            item: this.props.item,
        });
    }

    onLoad = (reactFlowInstance) => {
        reactFlowInstance.fitView();
        this.setState({reactFlowInstance});
    };

    render() {
      const { isLoggedIn, message, user } = this.props;
      const { elements, isInProgress, loading } = this.state;

      return (
        elements !== null &&
        <ReactFlow
          elements={elements}
          nodeTypes={nodeTypes}
          onLoad={this.onLoad}
        >
            <Background color="#aaa" gap={4} />
        </ReactFlow>
      );
  }
};

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    const { message } = state.message;
    return {
        isLoggedIn,
        message,
        user,
    };
}

export default connect(mapStateToProps)(MiniPipeline);
