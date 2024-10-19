import React from 'react';

import Button from "react-bootstrap/Button";
import Modal from "react-bootstrap/Modal";

class ConfirmationModal extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            show: this.props.show,
            onHide: this.props.onHide,
            title: this.props.title,
            body: this.props.body,
            onCancel: this.props.onCancel,
            onConfirm: this.props.onConfirm,
            cancelText: this.props.cancelText,
            confirmText: this.props.confirmText,
            darkTheme: localStorage.getItem('darkTheme') !== "false",
        };
    }

    UNSAFE_componentWillReceiveProps(nextProps, nextContext) {
        this.setState({
            show: this.props.show,
            onHide: this.props.onHide,
            title: this.props.title,
            body: this.props.body,
            onCancel: this.props.onCancel,
            onConfirm: this.props.onConfirm,
            cancelText: this.props.cancelText,
            confirmText: this.props.confirmText,
            darkTheme: localStorage.getItem('darkTheme') !== "false",
        });
    }

    render() {
        const { show, onHide, title, body, onCancel, onConfirm, cancelText, confirmText } = this.props;

        const { darkTheme } = this.state;

        return (
            <>
                <Modal show={show} onHide={onHide}>
                    <Modal.Header 
                        closeButton
                        closeVariant={
                            darkTheme ? 'white' : 'dark'
                        }
                    >
                        <Modal.Title>{title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {body}
                    </Modal.Body>
                    <Modal.Footer>
                        <Button variant="light" onClick={onCancel}>
                            {cancelText ?? "Cancel"}
                        </Button>
                        <Button variant="primary" onClick={onConfirm}>
                            {confirmText}
                        </Button>
                    </Modal.Footer>
                </Modal>
            </>
        );
    }
}

export default ConfirmationModal;
