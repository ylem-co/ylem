import React from 'react';

import Modal from "react-bootstrap/Modal";

class RightModal extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            show: this.props.show,
            onHide: this.props.onHide,
            title: this.props.title,
            content: this.props.content,
            showHideButton: this.props.showHideButton,
            hidingPropery: this.props.hidingPropery,
            darkTheme: localStorage.getItem('darkTheme') !== "false",
        };
    }

    UNSAFE_componentWillReceiveProps(nextProps, nextContext) {
        this.setState({
            show: this.props.show,
            onHide: this.props.onHide,
            title: this.props.title,
            content: this.props.content,
            showHideButton: this.props.showHideButton,
            hidingPropery: this.props.hidingPropery,
            darkTheme: localStorage.getItem('darkTheme') !== "false",
        });
    }

    render() {
        const { show, onHide, title, content } = this.props;

        const { darkTheme } = this.state;
        
        return (
            <>
                {/* showHideButton === true &&
                    <div className="hideArrowButton">
                        <Tooltip title="Close output log" placement="top">
                            <div className="ht">
                                &raquo;
                            </div>
                        </Tooltip>
                    </div>
                */}
                <Modal show={show} onHide={onHide} className="right">
                    <Modal.Header 
                        closeButton 
                        closeVariant={
                            darkTheme ? 'white' : 'dark'
                        }
                    >
                        <Modal.Title>{title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {content}
                    </Modal.Body>
                </Modal>
            </>
        );
    }
}

export default RightModal;
