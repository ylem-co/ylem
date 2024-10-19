import React from "react";

import Input from "./input.component";

class InputChips extends React.Component {
  state = {
    items: this.props.items || [],
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
    var values = paste.match(/[\w\d\.-]+/gi);

    if (values) {
      var toBeAdded = values.filter(value => !this.isInList(value));

      var items = [...this.state.items, ...toBeAdded];

      this.setState({items});

      this.handleItemsChange(items);
    }
  };

  isValid(value) {
    let error = null;

    if (this.isInList(value)) {
      error = `${value} has already been added.`;
    }

    if (!this.isValidInput(value)) {
      error = `${value} is not valid.`;
    }

    if (error) {
      this.setState({ error });

      return false;
    }

    return true;
  }

  isInList(value) {
    return this.state.items.includes(value);
  }

  isValidInput(value) {
    return /[\w\d\.-]+/.test(value);
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

        <Input
          className={"input " + (this.state.error ? " has-error" : undefined)}
          value={this.state.value}
          placeholder="Type or paste value and press `Enter`, `Tab` or `,` to add to the list ..."
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

export default InputChips;
