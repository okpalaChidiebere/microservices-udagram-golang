// This file can be replaced during build by using the `fileReplacements` array.
// `ng build --prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  appName: 'Udagram',
  apiHost: 'http://abeb639b4d92a45f7890cf74244cb757-1351125667.ca-central-1.elb.amazonaws.com:8080/api/v0' //Your frontend interacts with the reverseproxy through the loadbalacer which is more safe to the public. You cant access the internal cluster ipAddress if you are not in the Cluster!
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.
