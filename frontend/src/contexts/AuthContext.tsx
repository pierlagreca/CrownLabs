import { UserManager } from 'oidc-client';
import { createContext, FC, PropsWithChildren, useEffect, useRef } from 'react';
import {
  REACT_APP_CROWNLABS_OIDC_CLIENT_SECRET,
  REACT_APP_CROWNLABS_OIDC_PROVIDER_URL,
  REACT_APP_CROWNLABS_OIDC_CLIENT_ID,
} from '../env';
interface IAuthContext {
  isLoggedIn: boolean;
}

export const AuthContext = createContext<IAuthContext>({
  isLoggedIn: false,
});

const AuthContextProvider: FC<PropsWithChildren<{}>> = props => {
  const { children } = props;
  const manager = useRef<UserManager>();
  // setup event listeners for browser media-query API prefer-color-schema change
  // probably excessive but it doesn't hurt
  useEffect(() => {
    // console.log(
    //   'REACT_APP_CROWNLABS_OIDC_CLIENT_SECRET',
    //   REACT_APP_CROWNLABS_OIDC_CLIENT_SECRET
    // );
    // console.log(
    //   'REACT_APP_CROWNLABS_OIDC_PROVIDER_URL',
    //   REACT_APP_CROWNLABS_OIDC_PROVIDER_URL
    // );
    // console.log(
    //   'REACT_APP_CROWNLABS_OIDC_CLIENT_ID',
    //   REACT_APP_CROWNLABS_OIDC_CLIENT_ID
    // );
    // console.log(
    //   'REACT_APP_CROWNLABS_OIDC_REDIRECT_URI',
    //   REACT_APP_CROWNLABS_OIDC_REDIRECT_URI
    // );
    const redirectUri = window.location.href;
    manager.current = new UserManager({
      automaticSilentRenew: true,
      response_type: 'code',
      filterProtocolClaims: true,
      scope: 'openid ',
      loadUserInfo: true,
      client_secret: REACT_APP_CROWNLABS_OIDC_CLIENT_SECRET,
      authority: REACT_APP_CROWNLABS_OIDC_PROVIDER_URL,
      client_id: REACT_APP_CROWNLABS_OIDC_CLIENT_ID,
      redirect_uri: `${redirectUri}/authCallback`,
      post_logout_redirect_uri: `${redirectUri}/logout`,
    });
    manager.current
      .getUser()
      .then(res => {
        // console.log('user', res);
      })
      .catch(err => {
        console.error('ERR USE', err);
      });
    // // console.log('MANAGER', manager.current);
    // console.log('loc', window.location.href);
    // manager.current.signinRedirect();
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn: false }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContextProvider;
