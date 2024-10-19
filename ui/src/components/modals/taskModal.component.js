import React from 'react';

import Button from "react-bootstrap/Button";
import Modal from "react-bootstrap/Modal";

class TaskModal extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            show: this.props.show,
            onHide: this.props.onHide,
            title: this.props.title,
            content: this.props.content,
            darkTheme: localStorage.getItem('darkTheme') !== "false",
        };
    }

    UNSAFE_componentWillReceiveProps(nextProps, nextContext) {
        this.setState({
            show: this.props.show,
            onHide: this.props.onHide,
            title: this.props.title,
            content: this.props.content,
            darkTheme: localStorage.getItem('darkTheme') !== "false",
        });
    }

    render() {
        const { show, onHide, title, content, confirmButton, submitButtonLoading, confirmText, onConfirm } = this.props;

        const { darkTheme } = this.state;
        
        return (
            <>
                <Modal show={show} onHide={onHide} className="taskModal" backdropClassName="taskModalBg" centered size="xl">
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
                    {confirmButton === true &&
                    <Modal.Footer>
                        <Button variant="primary" onClick={onConfirm}>
                            {submitButtonLoading && (
                                <span className="spinner-border spinner-border-sm spinner-primary"></span>
                            )}
                            {confirmText}
                        </Button>  
                    </Modal.Footer>
                    } 
                </Modal>
            </>
        );
    }
}

export default TaskModal;
