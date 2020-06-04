'use strict';
exports.main_handler = async (event, context, callback) => {
  console.log(event)
  console.log(context)
  let crm_header = {
    "userId": parseInt(event.pathParameters.id)
  }

  return Buffer.from(JSON.stringify(crm_header)).toString('base64')
};
