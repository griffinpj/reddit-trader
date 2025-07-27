import { Layout, Button, Checkbox, Form, Input, Flex, Typography, theme, Card } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
const { Content, Header } = Layout;
const { Title } = Typography;

import Signin from './forms/Signin';
import Signup from './forms/Signup';


function Auth(props) {
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    let Form = null;
    switch (props.type) {
        case 'login':
            Form = <Signin />;
            break;
        case 'register':
            Form = <Signup />;
            break;
    }

    return (
        <Layout>
            <Header style={{ display: 'flex', alignItems: 'center' }}>
                <div className="demo-logo" />
            </Header>
            <Content>
                <div
                    style={{
                        background: colorBgContainer,
                        minHeight: 280,
                        padding: 24,
                        borderRadius: borderRadiusLG
                    }}
                >
                    <Flex
                        vertical
                        align="center"
                        justify='center'
                        style={{
                            padding: 24
                        }}
                    >
                    {Form}
                    </Flex>
                </div>
            </Content>
        </Layout>
    );
};

export default Auth;
