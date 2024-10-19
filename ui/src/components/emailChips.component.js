import React from "react";

class EmailChips extends React.Component {
  state = {
    items: [],
    value: "",
    error: null
  };

  handleItemsChange = (items) => {
    this.props.changesHandler(items);
  }

  handleKeyDown = evt => {
    if (["Enter", "Tab", ","].includes(evt.key)) {
      evt.preventDefault();

      var value = this.state.value.trim();

      if (value && this.isValid(value)) {
        var items = [...this.state.items, this.state.value];

        this.setState({
          items: items,
          value: ""
        });

        this.handleItemsChange(items);
      }
    }
  };

  handleChange = evt => {
    this.setState({
      value: evt.target.value,
      error: null
    });
  };

  handleDelete = item => {
    var items = this.state.items.filter(i => i !== item);

    this.setState({items});

    this.handleItemsChange(items);
  };

  handlePaste = evt => {
    evt.preventDefault();

    var paste = evt.clipboardData.getData("text");
    var emails = paste.match(/[\w\d\.-]+@[\w\d\.-]+\.[\w\d\.-]+/g);

    if (emails) {
      var toBeAdded = emails.filter(email => !this.isInList(email));

      var items = [...this.state.items, ...toBeAdded];

      this.setState({items});

      this.handleItemsChange(items);
    }
  };

  isValid(email) {
    let error = null;

    if (this.isInList(email)) {
      error = `${email} has already been added.`;
    }

    if (!this.isEmail(email)) {
      error = `${email} is not a valid email address.`;
    }

    if (error) {
      this.setState({ error });

      return false;
    }

    return true;
  }

  isInList(email) {
    return this.state.items.includes(email);
  }

  isEmail(email) {
    return /[\w\d\.-]+@[\w\d\.-]+\.[\w\d\.-]+/.test(email);
  }

  render() {
    return (
      <>
        <div className="chips">
        {this.state.items.map(item => (
          <div className="tag-item" key={item}>
            {item}
            <button
              type="button"
              className="button"
              onClick={() => this.handleDelete(item)}
            >
              &times;
            </button>
          </div>
        ))}

        <input
          className={"input " + (this.state.error ? " has-error" : undefined)}
          value={this.state.value}
          placeholder="Type or paste email addresses and press `Enter`, `Tab` or `,` to add to the list ..."
          onKeyDown={this.handleKeyDown}
          onChange={this.handleChange}
          onPaste={this.handlePaste}
        />

        {this.state.error && <p className="error">{this.state.error}</p>}
        </div>
      </>
    );
  }
}

export default EmailChips;
