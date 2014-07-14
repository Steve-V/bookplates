#!/usr/bin/perl -w

# Converts QuickMark-style barcode scanner log to curl-able HTTP
# POST bodies for BookCrossing.

use strict;

use Text::CSV;
use WWW::Scraper::ISBN;
use URI::Escape;

my $csv = new Text::CSV;
my @lines = reverse <>;
my $scraper = new WWW::Scraper::ISBN;
$scraper->drivers("GoogleBooks", "Yahoo", "LOC", "ISBNnu");

my $bcid = undef;
my $result = undef;

# Known elephant in Cairo
push @lines, "qr,";

foreach my $line (@lines) {
	$line =~ s/;$//;
	unless ($csv->parse($line)) {
		warn "Unparseable ", $line, "\n";
		next;
	}
	my ($type, $code) = $csv->fields();
	if ($type eq "qr") {
		if ($result and $bcid) {
			# upload!
			warn "Uploading ", $result->book->{'title'}, "\n";
			my %fields = (
				Title => $result->book->{'title'},
				Author => $result->book->{'author'},
				ISBN => $result->isbn,
				CategoryId => 0,
				Status => 'PermanentCollection',
				Comments => 'On the Godspeed bookshelf.',
				BCID => $bcid,
				Asin => Business::ISBN->new($result->isbn)->as_isbn10->as_string([]),
			);
			my @kvpairs;
			while(my ($k, $v) = each %fields) {
				push @kvpairs, join("=", $k, uri_escape_utf8($v));
			}
			print join("&", @kvpairs), "\n";
		}
		$result = undef;
		if ($code =~ m|/([-\d+]+)$|) {
			$bcid = $1;
			warn "New BCID ", $bcid;
		} else {
			warn "Unparseable BCID ", $code;
		}
	} else {
		if (!$result) {
			unless (Business::ISBN->new($code)->is_valid) {
				warn "Invalid ISBN $code";
				next;
			}
			my $r = $scraper->search($code);
			if ($r->found) {
				$result = $r;
				warn "Found ", $code, "\n";
			} else {
				warn "Not found ", $code, "\n";
			}
		}
	}
}
