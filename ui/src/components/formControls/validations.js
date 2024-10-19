export const required = (value, props) => {
    if (!value|| (props.isCheckable && !props.checked)) {
        return (
            <div role="alert" className="invalidInput"></div>
        );
    }
};

export const integrationBody = (value, props) => {
    if (
        props.integration !== "jenkins"
        && (!value || (props.isCheckable && !props.checked))
    ) {
        return (
            <div role="alert" className="invalidInput"></div>
        );
    }
};

export const requiredForCreation = (value, props) => {
    if (!value && props.iscreation && props.iscreation === "true") {
        return (
            <div role="alert" className="invalidInput"></div>
        );
    }
};

export const isEqual = (value, props, components) => {
    const bothUsed = components.password[0].isUsed && components.confirmPassword[0].isUsed;
    const bothChanged = components.password[0].isChanged && components.confirmPassword[0].isChanged;

    if (bothChanged && bothUsed && components.password[0].value !== components.confirmPassword[0].value) {
        return <div role="alert" className="invalidInput"></div>;
    }
};

export const isCron = (value) => {
    //const regExp = new RegExp(/^(?#minute)(\*|(?:[0-9]|(?:[1-5][0-9]))(?:(?:\-[0-9]|\-(?:[1-5][0-9]))?|(?:\,(?:[0-9]|(?:[1-5][0-9])))*)) (?#hour)(\*|(?:[0-9]|1[0-9]|2[0-3])(?:(?:\-(?:[0-9]|1[0-9]|2[0-3]))?|(?:\,(?:[0-9]|1[0-9]|2[0-3]))*))(?#day_of_month)(\*|(?:[1-9]|(?:[12][0-9])|3[01])(?:(?:\-(?:[1-9]|(?:[12][0-9])|3[01]))?|(?:\,(?:[1-9]|(?:[12][0-9])|3[01]))*)) (?#month)(\*|(?:[1-9]|1[012]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?:(?:\-(?:[1-9]|1[012]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC))?|(?:\,(?:[1-9]|1[012]|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC))*)) (?#day_of_week)(\*|(?:[0-6]|SUN|MON|TUE|WED|THU|FRI|SAT)(?:(?:\-(?:[0-6]|SUN|MON|TUE|WED|THU|FRI|SAT))?|(?:\,(?:[0-6]|SUN|MON|TUE|WED|THU|FRI|SAT))*))$/);

    const regExp = new RegExp(/(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|([\d\*]+(\/|-)\d+)|\d+|\*) ?){5,7})/);
  
    if (value && !regExp.test(value)) {
        return (
            <div role="alert" className="invalidInput"></div>
        );
    }
}

export const isTrialHost = (value) => {
    const regExp = new RegExp(/(Trial Host|Demo Host)/i);
  
    if (value && regExp.test(value)) {
        return true;
    }

    return false;
}
