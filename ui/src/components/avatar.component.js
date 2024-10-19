import React from 'react';

import Gravatar from 'react-gravatar';

class Avatar extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            size: this.props.size,
            avatar_url: this.props.avatar_url,
            email: this.props.email
        };
    }

    UNSAFE_componentWillReceiveProps(nextProps, nextContext) {
        this.setState({
            size: this.props.size,
            avatar_url: this.props.avatar_url,
            email: this.props.email,
        });
    }

    render() {
        var avatar_url = this.props.avatar_url;
        var size = this.props.size;
        var email = this.props.email;
        var circled = this.props.circled;

        return (
            <>
                <div className={circled ? "c-avatar circledAvatar" : "c-avatar squaredAvatar"}>
                    {circled &&
                        <div
                            className="circledAvatarBorder"
                            style={{
                                "width": size+22,
                                "height": size+22,
                                "borderRadius": (size+22)/2,
                                "borderColor": "#fff"
                            }}
                        ></div>
                    }
                {
                    avatar_url !== null && avatar_url !== ""
                        ? <img src={avatar_url} alt="" width={size} className="dtmnAvatar"/>
                        : <Gravatar
                            email={email !== null && email !== "" ? email : 'support@ylem.co'}
                            size={size}
                            default="mp" />
                }
                </div>
            </>
        );
    }
}

export default Avatar;
