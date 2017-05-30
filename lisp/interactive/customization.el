;; {{ define "public/customization" }}
;; Package related configuration.

(defgroup Go.el nil
  "Configure Go->Emacs Lisp translator."
  :group 'development)

(defcustom Go-utils-path ""
  "Specifies path which is used to find required binaries.
If remain empty, system PATH is used."
  :group 'Go.el
  :type 'directory)

(defcustom Go-emacs-pkg-root ""
  "Path to `emacs' Go package."
  :group 'Go.el
  :type 'directory)

;; {{ end }}