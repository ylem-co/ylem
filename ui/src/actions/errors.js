export const HTTP_ERRORS = {
    400: {
        'email': 'Invalid e-mail address\n',
        'phone': 'Invalid phone number\n',
        'number': 'Invalid phone number\n',
        'name': 'Invalid name\n',
        'host': 'Invalid host\n',
        'severity': 'Invalid severity\n',
        'sql_query': 'SQL query should be valid and must not contain DROP, TRUNCATE, or DELETE commands\n',
        'source_uuid': 'Invalid source\n',
        'expression': 'Invalid expression\n',
        'code': 'Invalid confirmation code\n',
        'ssh_port': 'Invalid SSH port\n',
        'port': 'Invalid database port\n',
        'user': 'Invalid user\n',
        'credentials': 'Invalid credentials\n',
        'project_id': 'Invalid project ID\n',
        'headers': 'Headers must contain a valid JSON\n',
        'payload': 'Payload must contain a valid JSON\n',
        'database': 'Invalid database\n',
        'channel': 'Invalid Slack channel\n',
        'authorization_uuid': 'Invalid Slack authorization\n',
        'type': 'Invalid type\n',
        'connection_type': 'Invalid connection type\n',
        'organization_name': 'Invalid organization Name\n',
        'merge.field_names':  'Invalid list of fields\n',
        'organization_exists': 'Organization already exists\n',
        'env_variable_name': 'Name is invalid. It can only contain digits, letters and symbols - and _\n',
        'env_variable_value': 'Value is invalid. It can only contain integers, decimals, and strings with letters and symbols - and _\n',
        'env_variable_name_exists': 'Such environment variable already exists',
        'folder_name_exists': 'Such folder already exists',
        'run_pipeline_uuid': 'A pipeline cannot be triggered from itself',
        'run_pipeline_uuid_run': 'You cannot run this pipeline. Please make sure you have correct access rights to it',
        'password': 'Password must be at least 8 symbols long and must contain:\n' +
            '- At least one digit\n' +
            '- At least one lowercase character\n' +
            '- At least one uppercase character\n' +
            '- At least one special character\n',
        'conirm_password': 'Password and its confirmation do not match\n',
        'Condition.expression: Missing required field': 'Expression cannot be empty\n',
        'Aggregator.expression: Missing required field': 'Expression cannot be empty\n',
    },
};

export const prepareErrorMessage = (error) => {
    let message = '';
    if (
        error.response.status &&
        error.response.status in HTTP_ERRORS &&
        error.response.data.fields
    ) {
        let errorFields = error.response.data.fields.split(",");
        for (let key in HTTP_ERRORS[error.response.status]) {
            if (
                errorFields.includes(key)
            ) {
                message += HTTP_ERRORS[error.response.status][key] + '\n';
            }
        }

        if (message === '') {
            message = error.response.data.fields;
        }
    } else {
        message =
            (error.response &&
                error.response.data &&
                error.response.data.message) ||
            error.message ||
            error.toString();
    }

    message = message.trim();

    return message;        
}
