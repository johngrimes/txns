var handleCatChangeError = function(xhr, status, errorThrown) {
  console.log(xhr);
  console.log(status);
  console.log(errorThrown);
  var alertEl = $('<div class="alert alert-warning" role="alert">There was an error saving your category change.</div>');
  var closeButton = $('<button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>');
  $(alertEl).append(closeButton);
  $('body > .container:first-child').before(alertEl);
  return true;
};

$(document).ready(function() {
  $('.txns-category-select').change(function() {
    var form = $(this).parents('form').first();
    var action = $(form).attr('action');
    var value = this.value == 0 ? null : this.value;
    if (this.value == 'new-category') {
      $('#modal-new-category').modal();
      $('#modal-new-category').on('shown.bs.modal', function (e) {
        $('#new-category-name').focus();
      });
      return false;
    }
    $.ajax({
      url: action,
      data: JSON.stringify({
        'txn': {
          'category_id': value
        }
      }),
      type: 'PATCH',
      error: handleCatChangeError
    });
    if (this.value == '0') {
      $(this).parents('tr').first().addClass('uncategorised');
    } else {
      $(this).parents('tr').first().removeClass('uncategorised');
    }
    return true;
  });
});
