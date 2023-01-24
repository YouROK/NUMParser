$(document).ready(function() {

    // scroll to top
    if ($('#back-to-top').length) {
        var scrollTrigger = 100, // px
            backToTop = function () {
                var scrollTop = $(window).scrollTop();
                if (scrollTop > scrollTrigger) {
                    $('#back-to-top').addClass('show');
                } else {
                    $('#back-to-top').removeClass('show');
                }
            };
        backToTop();
        $(window).on('scroll', function () {
            backToTop();
        });
        $('#back-to-top').on('click', function (e) {
            e.preventDefault();
            $('html,body').animate({
                scrollTop: 0
            }, 700);
        });
    }
    //

    $('#clear').on('click', function (e) {
        $('#s').val('');
    });

    $('input[type=radio]').change(function () {
            if($(this).is(":checked")) {
                $('form#searchForm').trigger("submit");
            }
        }
    );

    $('#proxyTrackers').on('click', function (e) {
        if($(this).is(":checked")) {
            $.cookie("proxy_trackers", 1);
        } else  {
            $.cookie("proxy_trackers", 0);
        }
    });

    // * React search!
    // * @floatrx, 24Jan/2019
    // * modified @Paukan

    $(function() {
        var $searchForm = $('form#searchForm'),
            $results = $('#resultsDiv'),
            $s = $searchForm.find('#s');

        if ($searchForm.length && $results.length) {
            // $('body').addClass('is-search')
            function updateSearchResults(data) {
                // Create virtual DOM from recieved data – $(data)
                // Cause: This is easier way to parse responsed page
                var $recievedContent = $(data).find('#resultsDiv');
                // Element exist
                if ($recievedContent.length) {
                    $results.html($recievedContent.html())
                } else {
                    console.log('Error. Update search results')
                }
            }
            console.log('Catch default submit action!');

            $searchForm.submit(function(e) {
                var queryString = $s.val(),
                    sortSize  = $searchForm.find('#sortSize').is(':checked'),
                    sortDate = $searchForm.find('#sortDate').is(':checked'),
                    forcedSearch = $searchForm.find('#forcedSearch').is(':checked'),
                    // proxyTrackers = $searchForm.find('#proxyTrackers').is(':checked'),
                    timerId = 'Found "' + queryString + '" by (ms)';

                if (queryString.length > 0) {

                    $.get({
                        url: '/' + queryString + (sortSize ? '/size' : (sortDate ? '/date' : '')) + (forcedSearch ? '?forced=1' : ''), // + (proxyTrackers ? '&pxy=1' : ''),
                        beforeSend: function () {
                            $results.disable(true); // lock before request
                            console.time(timerId)
                        }
                    }).done(function (data, jqXHR, statusText) {
                        // Update histofy state and document titls
                        document.title = 'Скачать торрент - ' + queryString + ' | Торренты в BDRip и HDRip качестве';
                        history.replaceState({}, document.title, location.origin + '/' + queryString);
                        // Update search results without page reloading
                        // No JSON!
                        updateSearchResults(data)
                    }).fail(function (data, statusText) {

                        // update even if not found (show "not found" message)
                        if (data.status == 404) {
                            updateSearchResults(data.responseText);
                        }

                        console.error('Fail. ', data, statusText) // console log fails
                    }).always(function () {
                        $results.disable(); // unlock results
                        console.timeEnd(timerId)
                    });
                } else {
                    $searchForm.find('#validateError').html("Поисковый запрос пустой!");

                }
                return false // disable form submitting
            });

            // remove error msg
            $s.on("input", function () {
                if ($s.val().length > 0) {
                    $searchForm.find('#validateError').html("");
                }
            });

        }
    });

}); // document.ready

// Extend jQuery
$.fn.disable = function(bool) {
    var $this = $(this);
    if (!bool) {
        $this.removeAttr('disabled') } else {
        $this.attr('disabled', 'true') }
    return $this
}