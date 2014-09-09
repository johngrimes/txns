/** @jsx React.DOM */

var AccountPage = React.createClass({
  handleSubmitSuccess: function() {
    this.refs.txnTable.refreshTxns();
    this.refs.txnUploadForm.render();
  },

  render: function() {
    return (
      <section>
        <button type="button" className="btn btn-primary btn-upload-statement">Upload a statement</button>
        <TxnTable ref="txnTable" />
        <TxnUploadForm ref="txnUploadForm" onSubmitSuccess={this.handleSubmitSuccess} />
      </section>
    );
  }
});

var TxnTable = React.createClass({
  getInitialState: function() {
    return { txns: [] };
  },

  getDefaultProps: function() {
    return {
      txnsURL: 'http://localhost:4000/txns',
      pollFrequency: 10000
    };
  },

  componentDidMount: function() {
    this.refreshTxns();
    setInterval(this.refreshTxns, this.props.pollFrequency);
  },

  refreshTxns: function() {
    $.getJSON(this.props.txnsURL, function(data) {
      this.setState({ txns: data });
    }.bind(this));
  },

  render: function() {
    var txns = this.state.txns;
    return (
      <table className="table table-bordered table-hover">
        <thead>
          <tr>
            <th>Date</th>
            <th>Description</th>
            <th>Debit</th>
            <th>Credit</th>
            <th>Balance</th>
          </tr>
        </thead>
        <tbody>
          {txns.map(function(txn) {
            return <Txn key={txn["id"]} txn={txn} />;
          })}
        </tbody>
      </table>
    );
  }
});

var Txn = React.createClass({
  render: function() {
    var txn = this.props.txn;
    return (
      <tr>
        <td>{txn.date}</td>
        <td>{txn.description}</td>
        <td>{txn.debitCents}</td>
        <td>{txn.creditCents}</td>
        <td>{txn.balanceCents}</td>
      </tr>
    );
  }
});

var TxnUploadForm = React.createClass({
  getDefaultProps: function() {
    return {
      action: 'http://localhost:4000/txns'
    };
  },

  handleSubmit: function() {
    formData = new FormData();
    fileInput = this.refs.fileInput.getDOMNode();
    formData.append(fileInput.name, fileInput.files[0]);
    $.ajax({
      url: this.props.action,
      type: 'POST',
      contentType: false,
      processData: false,
      data: formData,
      success: function(data) {
        this.props.onSubmitSuccess();
      }.bind(this),
      error: function(xhr, err) {
        //TODO: Give error message back to user.
        console.log(err);
      }
    });
    return false;
  },

  render: function() {
    return (
      <div id="modal-upload-statement" className="modal fade">
        <div className="modal-dialog">
          <div className="modal-content">
            <form onSubmit={this.handleSubmit} action={this.props.action} method="POST" encType="multipart/form-data">
              <div className="modal-header">
                <button type="button" className="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span className="sr-only">Close</span></button>
                <h4 className="modal-title">Upload a statement</h4>
              </div>
              <div className="modal-body">
                <input type="file" ref="fileInput" name="txn_file" />
              </div>
              <div className="modal-footer">
                <button type="submit" className="btn btn-primary">Upload</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
});

React.renderComponent(
  <AccountPage />,
  document.getElementById('account')
);
