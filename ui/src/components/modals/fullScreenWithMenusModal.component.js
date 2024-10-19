import React from 'react';

import Button from "react-bootstrap/Button";
import Modal from "react-bootstrap/Modal";

class FullScreenWithMenusModal extends React.Component {
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
        const { 
            show, 
            onHide, 
            title, 
            content, 
            confirmButton, 
            altButton, 
            onAltButtonClick, 
            altButtonText,
            altButton2, 
            onAltButton2Click, 
            altButton2Text,
        } = this.props;

        const { darkTheme } = this.state;

        return (
            <>
                <Modal show={show} onHide={onHide} fullscreen backdrop={false}>
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
                    {(altButton === true || altButton2 === true || confirmButton === true) &&
                    <Modal.Footer>
                        {altButton === true &&
                        <Button variant="light" onClick={onAltButtonClick}>
                            {altButtonText}
                        </Button>
                        }
                        {altButton2 === true &&
                        <Button variant="secondary" onClick={onAltButton2Click}>
                            {altButton2Text}
                        </Button>
                        }
                        {confirmButton === true &&
                        <Button variant="primary">
                            Save
                        </Button>
                        }
                    </Modal.Footer>
                    }
                </Modal>
            </>
        );
    }
}

export default FullScreenWithMenusModal;
