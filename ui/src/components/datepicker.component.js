import React from 'react';
import moment from 'moment';

import 'react-dates/initialize';
import 'react-dates/lib/css/_datepicker.css';
import { DateRangePicker } from 'react-dates';

class Datepicker extends React.Component {
    constructor(props) {
        super(props);

        let minDate = new Date(this.props.createdAt);
        minDate.setFullYear(minDate.getFullYear() - 1);

        this.state = {
            showDatePicker: this.props.showDatePicker,
            minDate: minDate,
            maxDate: this.props.maxDate ? new Date(this.props.maxDate) : null,
            endDate: this.props.dateTo,
            startDate: this.props.dateFrom,
        };
    }

    handleDateChange = (startDate, endDate) => {
        if (
            typeof startDate ===  'object'
            && startDate !== null
            && typeof endDate ===  'object'
            && endDate !== null
        ) {
            let sD = startDate.format("YYYY-MM-DD")
            let eD = endDate.format("YYYY-MM-DD");
            this.setState({ 
                startDate: sD, 
                endDate: eD 
            });

            this.props.changeDates(sD, eD);
        }
    };

    /*UNSAFE_componentWillReceiveProps(props) {
        let minDate = new Date(this.props.createdAt);
        minDate.setFullYear(minDate.getFullYear() - 1);

        this.setState({
            minDate: minDate,
            endDate: props.dateTo,
            startDate: props.dateFrom,
        });

        if (this.props.maxDate) {
            this.setState({ maxDate: this.props.maxDate });
        }

        this.setState({ showDatePicker: this.props.showDatePicker });
    }*/

    render() {
        const className = this.state.showDatePicker === true ? 'datePickerShow' : 'datePickerHide';

        return (
            <>
                <div className={className}>
                    <DateRangePicker
                        startDate={moment(this.state.startDate)} // momentPropTypes.momentObj or null,
                        startDateId="your_unique_start_date_id" // PropTypes.string.isRequired,
                        endDate={moment(this.state.endDate)} // momentPropTypes.momentObj or null,
                        endDateId="your_unique_end_date_id" // PropTypes.string.isRequired,
                        onDatesChange={({ startDate, endDate }) => this.handleDateChange(startDate, endDate)} // PropTypes.func.isRequired,
                        focusedInput={this.state.focusedInput} // PropTypes.oneOf([START_DATE, END_DATE]) or null,
                        onFocusChange={focusedInput => this.setState({ focusedInput })} // PropTypes.func.isRequired,
                        isOutsideRange={date => date.isBefore(this.state.minDate, 'day') || date.isAfter(this.state.maxDate !== null ? this.state.maxDate : new Date(), 'day')}
                        minimumNights={0}
                    />
                </div>
            </>
        );
    }
}

export default Datepicker
