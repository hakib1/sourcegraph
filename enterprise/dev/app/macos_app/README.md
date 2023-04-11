# putting together a macOS app bundle for Sourcegraph App

# setup

## Signing certificates

In order to distribute the app bundle, it needs to be code-signed using Apple certificates.

There are two kinds of certificates needed:

1. Developer ID Application
1. Signs all of the code in the app bundle, and the app bundle itself.
1. Developer ID Installer
1. Won't be needed as long as the method of distribution is a simple .dmg archive. If a EULA is added to the .dmg, or a .pkg is created, it will need to be signed by the Developer ID Installer certificate

If they have expired and need to be re-generated, login to https://developer.apple.com as the Account Holder and create new ones.

### get the certificates into Google Secret Manager

Login to https://developer.apple.com as the Account Holder.

Download the certificate, which will download an unencrypted `.cer` file.

For security, we need to use encrypted certificates. To turn the unencrypted `.cer` file into an encrypted `.p12` file, either import into Keychain Access, then export (from the "login" keychain, not the "system" keychain), selecting as the File Format "Personal Information Exchange (.p12)", or use `openssl` in this two-step process

```
printf 'strong, randm password' >cert.key
openssl x509 -in <cert name>.cer -inform DER -out <cert name>.pem -outform PEM
openssl pkcs12 -export -out <cert name>.p12 -inkey cert.key -in <cert name>.pem
```

Store the `p12` file in the

## App-specific credentials for notarization

## Secrets

Secrets can be found in https://docs.google.com/document/d/1YDqrpxJhdfudlRsYTQMGeGiklmWRuSlpZ7Xvk9ABitM/edit?usp=sharing (requires Sourcegraph team credentials)

Secrets are also stored in 1Password, in the Apple Developer vault

The secrets have been stored in Google Secrets Manager, and are set up in the buildkite config so they are pulled into the CI pipeline.

### code-signing secrets

- APPLE_DEV_ID_APPLICATION_CERT - the encrypted code-signing certificate file (.p12)
- APPLE_DEV_ID_APPLICATION_PASSWORD - password to the code-signing certificate file

### notarization secrets

- APPLE_APP_STORE_CONNECT_API_KEY_ID - App Store Connect API ID
- APPLE_APP_STORE_CONNECT_API_KEY_ISSUER - App Store Connect API Issuer GUID
- APPLE_APP_STORE_CONNECT_API_KEY_FILE - App Store Connect API key file (.p8)

## dependencies

### referenced in

- Xcode project - the `build_app_bundle_template.sh` shell script is set up to run after an Archive action - it downloads the dependencies in order to build the app bundle template

### git

We include `git` in the app bundle to avoid an external runtime dependeny on `git`.

We have built universal binaries from the latest version at the time - 2.39.2.

New universal binaries can be built using `enterprise/dev/app/macos_app/build_git_macos.sh`.

It has to be run on macOS, with Xcode installed.

Run `build_git_macos.sh --help` to see the options.

It uses default versions: gettext 0.21.1 and git 2.39.2 - pass other versions as options to the script if you need different versions.

The output is a gzipped tar archive in the working directory
(defaults to the current working directory; can be modified by using the `--workdir` option)
named with the format `git-universal-${VERSION}.tar.gz`

To store the archive where the macOS app bundle build process can get to it,
upload it to the `sourcegraph_app_macos_dependencies` GCS bucket:

```
gsutil cp git-universal-${VERSION}.tar.gz gs://sourcegraph_app_macos_dependencies
```

### src-cli

We include `src` in the app bundle so that it doesn't have to download it from elsewhere, like Homebrew.

New universal binaries can be built using `enterprise/dev/app/macos_app/build_src-cli_macos.sh`.
It can be run on Linux or macOS, or maybe even Windows. It will generate macOS universal binaries on any platform.
It has the option to either build from source, or download a release. Currently it downloads a release.

The output is a gzipped tar archive in the working directory
(defaults to the current working directory; can be modified by using the `--workdir` option)
named with the format `src-universal-${VERSION}.tar.gz`

To store the archive where the macOS app bundle build process can get to it,
upload it to the `sourcegraph_app_macos_dependencies` GCS bucket:

```
gsutil cp src-universal-${VERSION}.tar.gz gs://sourcegraph_app_macos_dependencies
```

### universal-ctags

We include a custom build of `universal-ctags` in the app bundle that is built with json support.

New universal binaries can be built using `enterprise/dev/app/macos_app/build_universal-ctags_macos.sh`.
It can be run on Linux or macOS, or maybe even Windows. It will generate macOS universal binaries on any platform.
It has the option to either build from source, or download a release. Currently it downloads a release.

The output is a gzipped tar archive in the working directory
(defaults to the current working directory; can be modified by using the `--workdir` option)
named with the format `universal-ctags-universal-${VERSION}.tar.gz`

To store the archive where the macOS app bundle build process can get to it,
upload it to the `sourcegraph_app_macos_dependencies` GCS bucket:

```
gsutil cp universal-ctags-universal-${VERSION}.tar.gz gs://sourcegraph_app_macos_dependencies
```

### dependencies

- autoconf
- autoreconf
- maybe other build tools

## build and deploy the app bundle template

The app bundle template is an Xcode project in the `app bundle Xcode project` directory. It packages dependencies and includes a binary that provides a management GUI and launches the `sourcegraph_launcher.sh` shell script.

To generate and upload to GCS a new app bundle template, click on **Product** in the menu bar, then **Archive** in the list. Xcode will build the project, archive it, and run the `build_app_bundle_template.sh` shell script that's part of the Xcode project. That shell script will download the aforementioned dependencies from the `sourcegraph_app_macos_dependencies` GCS bucket, extract and place them into the app bundle template, create a gzipped tar archive of the app bundle template, upload it to `sourcegraph_app_macos_dependencies` GCS bucket, named with the date and time of the **Archive** build, and create a `template-version.txt` file in the GCS bucket whose contents are the date and time, as a pointer to the new template archive.

If you don't have write permission to the `sourcegraph_app_macos_dependencies` GCS bucket, the gzipped tar archive will be placed in your Downloads directory, and you can arrange to get it uploaded from there. Be sure to create/update the `template-version.txt` file as well.

To obtain write access to the GCP bucket, follow [these instructions](https://handbook.sourcegraph.com/departments/security/tooling/entitle_request/) to request Storage Object Admin access to the `sourcegraph_app_macos_dependencies` GCS bucket, which is in the "Sourcegraph CI" project.

Values for the request form:

- Request type: Specific Permission
- Integration: GCP Development Projects
- Resource Types: buckets
- Resource: Sourcegraph CI/sourcegraph_app_macos_dependencies
- Role: Storage Object Admin
- Grant Method: Direct
- Permission Duration: 3 Hours
- Add justification: Upload new artifacts for the Sourcegraph App macOS app bundle

`build_app_bundle_template.sh` will write output to the log file `${HOME}/Downloads/build_app_bundle_template.log` as well as to Xcode logs.

## tools

### sourcegraph/apple-codesign

Signing macOS artifacts (executables and app bundles) is done using the [sourcegraph/apple-codesign](https://github.com/sourcegraph/apple-codesign) Docker image. If that image is not published to Docker Hub, it can be created locally. The shell script `setup_codesign.sh` does that. The current version of the Docker image is 0.22.0.

# Process macOS artifacts

There are several things that need to be done for the macOS artifacts

- sign and notarize the standalone executable
- build, sign and notarize the app bundle
- build and notarize the dmg

The shell script `post-process_macos_artifacts.sh` handles all of those steps.

## Sign and notarize the standalone executable

The executable is downloaded from the versioned GCS bucket, signed using `sign_macos_artifact.sh`, notarized using `notarize_macos_artifact.sh`, and uploaded back to the versioned GCS bucket. The `checksums.txt` file in the same bucket is updated with the new checksum of the signed executable.

## Build, sign and notarize the app bundle

The standalone executable from the previous step is passed to `build_macos_app.sh`, then `sign_macos_app.sh` is used to sign the app bundle, and `notarize_macos_artifact.sh` to notarize it. The app bundle is compressed into a zip archive and uploaded to the versioned GCS bucket.

## Build and notarize the dmg

If the `post-process_macos_artifacts.sh` script is running on macOS, it will build and notarize a dmg also. The dmg is not built outside of macOS because the process relies on macOS-specific tools. Linux-based tools have been tested to try to remove the macOS dependency, but none have been satisfactory. The process of building and notarizing a dmg is accomplished by first running `create_sourcegraph_app_dmg.sh` to create the dmg, then signing it using `sign_macos_artifact.sh`, notarizing it using `notarize_macos_artifact.sh`, and uploading it to the versioned GCS bucket.

# Scripts

## sign_macos_artifact.sh

Uses Apple code-signing certificates and the `sourcegraph/apple-codesign` Docker image to sign artifacts. Does not recurse into directories.

To maintain compatibility with `goreleaser`, accepts a file path to the artifact in the environemnt variable `artifact` and will place the signed artifact in the location specified by the environment variable `signature`, if it is a valid file path.

Also accepts a file path to the artifact on the command line, which trumps the `artifact` environment variable.

## notarize_macos_artifact.sh

Uses Apple App Store Connect API keys and the `sourcegraph/apple-codesign` Docker image to notarize an artifact. If the artifact is an app bundle, will compress into a zip archive in order to submit for notarizing. Otherwise, will submit only files, not directories.

To maintain compatibility with `goreleaser`, accepts a file path to the artifact in the environemnt variable `artifact` and will place the signed artifact in the location specified by the environment variable `signature`, if it is a valid file path.

Also accepts a file path to the artifact on the command line, which trumps the `artifact` environment variable.

## build_macos_app.sh

Builds a macOS app bundle by combining the `sourcegraph` executable binary specified in the `artifact` environment, `app_bundle/sourcegraph_launcher.sh`, and the app bundle template into an app bundle.

To maintain compatibility with `goreleaser`, accepts a file path to the `sourcegraph` executable binary in the environemnt variable `artifact` and will place the app bundle in the location specified by the environment variable `signature`, if it is a valid file path.

If `artifact` is not set, will download the executable binary from the GCS `sourcegraph-app-releases/${VERSION}/` bucket.

Will use the app bundle template specified by the `app_template_path` environment. If that environment variable is not set, will download the latest version of the app bundle template from GCS.

Names the app bundle according to the `app_name` environment variable, defaulting to "Sourcegraph App". If `signature` is set, will use that as the destination name, instead of `app_name`.

## sign_macos_app.sh

Uses Apple code-signing certificates and the `sourcegraph/apple-codesign` Docker image to sign a macOS app bundle.

To maintain compatibility with `goreleaser`, accepts a path to the app bundle in the environemnt variable `artifact` and will place the signed app bundle in the location specified by the environment variable `signature`, if it is a valid file path.

Also accepts a path to the app bundle on the command line, which trumps the `artifact` environment variable.

Unlike `sign_macos_artifact.sh`, `sign_macos_app.sh` will recurse into the app bundle directory, signing all macOS executables in it. It then signs the whole app bundle.

## create_sourcegraph_app_dmg.sh

Builds a dmg containing the macOS app. The is modified with a special background, resized to a particular dimension, and oriented vertically.

Requires macOS because it relies on `hdiutil` to create the volume and `osascript` to customize it.

The app bundle to wrap is passed as the first command line parameter; if missing it defaults to `${HOME}/Downloads/Sourcegraph App.app`.

# to-do

- strip and pack executables to reduce size
  - may not happen. upx-packed arm binaries won't run. upx-packed amd binaries run via Rosetta, but maybe not on Intel?